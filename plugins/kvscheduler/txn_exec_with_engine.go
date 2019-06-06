// Copyright (c) 2019 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kvscheduler

import (
	"runtime/trace"
	"time"

	"github.com/gogo/protobuf/proto"
	kvs "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	exec "github.com/ligato/vpp-agent/plugins/kvscheduler/internal/exec-engine"
	"github.com/ligato/vpp-agent/plugins/kvscheduler/internal/graph"
	"github.com/ligato/vpp-agent/plugins/kvscheduler/internal/utils"
)

// context associated with a running transaction
type txnContext struct {
	*transaction
	graphW   graph.RWAccess
	dryRun   bool
	executed kvs.RecordedTxnOps
}

// context associated with an in-progress key-value pair update
type kvContext struct {
	// input
	baseKey     string
	operation   kvs.TxnOperation
	metadata    kvs.Metadata
	origin      kvs.ValueOrigin

	// output
	newMetadata kvs.Metadata
	opRecord   *kvs.RecordedTxnOp

	// handlers
	node       graph.NodeRW
	descriptor *descriptorHandler

	// flags
	isDerived  bool
	isRecreate bool
	isRetry    bool

	// previous value state - for FinalizeTxnOperation
	prevValue   proto.Message
	prevUpdate  *LastUpdateFlag
	prevState   kvs.ValueState
	prevOp      kvs.TxnOperation
	prevErr     string
	prevDetails []string

	// for ExecuteTxnOperation
	dryExec   bool
	execStart time.Time
	execStop  time.Time
}

// executeTransactionWithEngine executes pre-processed transaction.
// If <dry-run> is enabled, Validate/Create/Delete/Update operations will not be executed
// and the graph will be returned to its original state at the end.
func (s *Scheduler) executeTransactionWithEngine(txn *transaction, graphW graph.RWAccess,
	dryRun bool) (executed kvs.RecordedTxnOps) {

	op := "execute transaction"
	if dryRun {
		op = "simulate transaction"
	}
	defer trace.StartRegion(txn.ctx, op).End()
	if dryRun {
		defer trackTransactionMethod("simulateTransaction")()
	} else {
		defer trackTransactionMethod("executeTransaction")()
	}

	txnCtx := &txnContext{
		transaction: txn,
		graphW:      graphW,
		dryRun:      dryRun,
	}

	var kvChanges []exec.KVChange
	for _, kv := range txn.values {
		kvChanges = append(kvChanges, exec.KVChange{
			Key:      kv.key,
			NewValue: kv.value,
			Context: &kvContext{
				baseKey:  kv.key,
				metadata: kv.metadata,
				origin:   kv.origin,
				isRetry:  txn.txnType == kvs.RetryFailedOps,
			},
		})
	}

	withRevert := txn.txnType == kvs.NBTransaction && txn.nb.revertOnFailure
	s.execEngine.RunTransaction(txnCtx, kvChanges, withRevert)

	// get rid of uninteresting intermediate pending Create/Delete operations
	executed = s.compressTxnOps(txnCtx.executed)
	return executed
}

// PrepareTxnOperation should update the underlying graph (abstracted-away
// at this level) to reflect the requested value change.
func (s *Scheduler) PrepareTxnOperation(txnPrivCtx exec.OpaqueCtx, kv *exec.KVChange, isRevert bool) (
	prevValue proto.Message) {

	txnCtx := txnPrivCtx.(*txnContext)
	kvCtx := kv.Context.(*kvContext)

	// obtain descriptor for the key
	descriptor := s.registry.GetDescriptorForKey(kv.Key)
	handler := newDescriptorHandler(descriptor)
	kvCtx.descriptor = handler

	// create new revision of the node for the given key-value pair
	kvCtx.node = txnCtx.graphW.SetNode(kv.Key)

	// remember previous value for a potential revert
	kvCtx.prevValue = kvCtx.node.GetValue()
	prevValue = kvCtx.prevValue

	// remember previous value status to detect and notify about changes
	kvCtx.prevUpdate = getNodeLastUpdate(kvCtx.node)
	kvCtx.prevState = getNodeState(kvCtx.node)
	kvCtx.prevOp = getNodeLastOperation(kvCtx.node)
	kvCtx.prevErr = getNodeErrorString(kvCtx.node)
	kvCtx.prevDetails = getValueDetails(kvCtx.node)

	// prepare operation description - fill attributes that we can even before
	// executing the operation
	kvCtx.opRecord = s.preRecordTxnOp2(kv, kvCtx, isRevert)

	// determine the operation type
	// TODO: take into consideration multiple sources
	//  - create/delete can actually result in Update
	if kv.KeepValue {
		kvCtx.operation = kvs.TxnOperation_UNDEFINED // determined in IsTxnOperationReady
	} else if kv.NewValue == nil {
		kvCtx.operation = kvs.TxnOperation_DELETE
	} else if kvCtx.node.GetValue() == nil || !isNodeAvailable(kvCtx.node) {
		kvCtx.operation = kvs.TxnOperation_CREATE
	} else {
		kvCtx.operation = kvs.TxnOperation_UPDATE
	}
	if txnCtx.dryRun {
		// do not actually execute the operation, just pretend
		kvCtx.dryExec = true
	}

	// update the lastUpdate flag
	lastUpdateFlag := &LastUpdateFlag{
		txnSeqNum: txnCtx.seqNum,
		txnOp:     kvCtx.operation,
		revert:    isRevert,
	}
	if kv.KeepValue {
		lastUpdateFlag.value = prevValue
	} else {
		lastUpdateFlag.value = kv.NewValue
	}
	if txnCtx.txnType == kvs.NBTransaction {
		lastUpdateFlag.retryEnabled = txnCtx.nb.retryEnabled
		lastUpdateFlag.retryArgs = txnCtx.nb.retryArgs
	} else if kvCtx.prevUpdate != nil {
		// inherit retry arguments from the last NB txn for this value
		lastUpdateFlag.retryEnabled = kvCtx.prevUpdate.retryEnabled
		lastUpdateFlag.retryArgs = kvCtx.prevUpdate.retryArgs
	} else if kvCtx.isDerived {
		// inherit from the parent value
		parentNode := txnCtx.graphW.GetNode(kvCtx.baseKey)
		prevParentUpdate := getNodeLastUpdate(parentNode)
		if prevParentUpdate != nil {
			lastUpdateFlag.retryEnabled = prevParentUpdate.retryEnabled
			lastUpdateFlag.retryArgs = prevParentUpdate.retryArgs
		}

	}
	kvCtx.node.SetFlags(lastUpdateFlag)

	// if the value is already "broken" by this transaction, do not try to update
	// anymore, unless this is a revert
	// (needs to be refreshed first in the post-processing stage)
	failingValue := kvCtx.prevState == kvs.ValueState_FAILED || kvCtx.prevState == kvs.ValueState_RETRYING
	updatedByThisTxn := kvCtx.prevUpdate != nil && kvCtx.prevUpdate.txnSeqNum == txnCtx.seqNum
	if failingValue && !isRevert && updatedByThisTxn {
		skipExec = true
		return
	}

	// prepare for the selected operation
	switch kvCtx.operation {
	case kvs.TxnOperation_DELETE:
		s.prepareForDelete(txnCtx, kv, kvCtx, isRevert)
	case kvs.TxnOperation_CREATE:
		skipExec = s.prepareForCreate(txnCtx, kv, kvCtx, isRevert)
	case kvs.TxnOperation_UPDATE:
		s.prepareForUpdate(txnCtx, kv, kvCtx, isRevert)
	}
	return
}

func (s *Scheduler) prepareForDelete(txnCtx *txnContext, kv *exec.KVChange, kvCtx *kvContext, isRevert bool) {

	// TODO
	return
}

func (s *Scheduler) prepareForCreate(txnCtx *txnContext, kv *exec.KVChange, kvCtx *kvContext, isRevert bool) (skipExec bool) {
	node := kvCtx.node

	// TODO

	// add descriptor flag
	if !kvCtx.descriptor.isNil() {
		node.SetFlags(&DescriptorFlag{kvCtx.descriptor.name()})
		node.SetLabel(kvCtx.descriptor.keyLabel(kv.Key))
	}

	// mark derived value
	if kvCtx.isDerived {
		kvCtx.node.SetFlags(&DerivedFlag{baseKey: kvCtx.baseKey})
	}

	// handle unimplemented value
	unimplemented := kvCtx.origin == kvs.FromNB && !kvCtx.isDerived && kvCtx.descriptor.isNil()
	if unimplemented {
		skipExec = true
		if getNodeState(kvCtx.node) == kvs.ValueState_UNIMPLEMENTED {
			// already known
			return
		}
		node.SetFlags(&UnavailValueFlag{})
		node.DelFlags(ErrorFlagIndex)
		kvCtx.opRecord.NOOP = true
		kvCtx.opRecord.NewState = kvs.ValueState_UNIMPLEMENTED
		s.updateNodeState(node, kvCtx.opRecord, args)
		return
	}

	// TODO: node.SetValue() should use merge of all values from all sources of this KV

	return
}

func (s *Scheduler) prepareForUpdate(txnCtx *txnContext, kv *exec.KVChange, kvCtx *kvContext, isRevert bool) {

	// TODO
	return
}

// IsTxnOperationReady should determine whether to:
//  - proceed with operation execution
//  - wait for some other key-value pairs (of the same transaction) to
//    be changed first - once those values are finalized, the readiness
//    check is repeated and the value change process continues accordingly
//  - block (freeze) some other values from entering the state machine while
//    this value is waiting/being executed (unfrozen when finalized)
func (s *Scheduler) IsTxnOperationReady(txnPrivCtx exec.OpaqueCtx, kv *exec.KVChange) (
	skipExec bool, precededBy []exec.KVChange, freeze utils.KeySet) {

	//txnCtx := txnPrivCtx.(*txnContext)
	kvCtx := kv.Context.(*kvContext)

	// TODO

	if kvCtx.operation == kvs.TxnOperation_UNDEFINED {
		kvCtx.operation = s.determineDepUpdateOperation(kvCtx.node)
		getNodeLastUpdate(kvCtx.node).txnOp = kvCtx.operation
	}

	return false, nil, nil
}

// ExecuteTxnOperation is run from another go routine by one of the workers.
// The method should apply the value change (call Create/Delete/Update on the
// associated descriptor) and return error if the operation failed.
func (s *Scheduler) ExecuteTxnOperation(workerID int, kv *exec.KVChange) (err error) {
	kvCtx := kv.Context.(*kvContext)
	node := kvCtx.node

	kvCtx.execStart = time.Now()
	defer func() { kvCtx.execStop = time.Now() }()

	if kvCtx.dryExec {
		return nil
	}
	switch kvCtx.operation {
	case kvs.TxnOperation_DELETE:
		err = kvCtx.descriptor.delete(node.GetKey(), node.GetValue(), node.GetMetadata())
	case kvs.TxnOperation_CREATE:
		kvCtx.newMetadata, err = kvCtx.descriptor.create(node.GetKey(), node.GetValue())
	case kvs.TxnOperation_UPDATE:
		kvCtx.newMetadata, err = kvCtx.descriptor.update(
			node.GetKey(), kvCtx.prevValue, node.GetValue(), node.GetMetadata())
	}
	return err
}

// FinalizeTxnOperation is run after the operation has been executed/skipped.
// Some more key-value pair may be requested to be changed as a consequence
// (followUp).
func (s *Scheduler) FinalizeTxnOperation(txnPrivCtx exec.OpaqueCtx, kv *exec.KVChange,
	wasRevert bool, opRetval error) (followUp []exec.KVChange) {

	// TODO - don't forget about all the skips (unimplemented, SB) and new metadata
	txnCtx := txnPrivCtx.(*txnContext)
	kvCtx := kv.Context.(*kvContext)

	// detect value state changes
	if !txnCtx.dryRun {
		nodeR := txnCtx.graphW.GetNode(kv.Key)
		if kvCtx.prevUpdate == nil || kvCtx.prevState != getNodeState(nodeR) || kvCtx.prevOp != getNodeLastOperation(nodeR) ||
			kvCtx.prevErr != getNodeErrorString(nodeR) || !equalValueDetails(kvCtx.prevDetails, getValueDetails(nodeR)) {
			s.updatedStates.Add(kvCtx.baseKey)
		}
	}

	return nil
}

// PrepareForTxnRevert is run before reverting of already applied key-value
// changes is started (due to error(s)).
func (s *Scheduler) PrepareForTxnRevert(txnPrivCtx exec.OpaqueCtx, failedKVChanges utils.KeySet) {
	txnCtx := txnPrivCtx.(*txnContext)
	s.refreshGraph(txnCtx.graphW, failedKVChanges, nil, true)
}
