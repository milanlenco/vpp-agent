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
	"sort"
	"time"

	"github.com/ligato/cn-infra/logging"
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
	newMergedValue proto.Message
	operation      kvs.TxnOperation

	// output
	newMetadata kvs.Metadata
	opRecord    *kvs.RecordedTxnOp

	// handlers
	node       graph.NodeRW
	descriptor *descriptorHandler

	// flags
	isRecreate bool
	isRetry    bool

	// previous value state - for FinalizeTxnOperation
	prevValue   proto.Message // merged
	prevUpdate  *LastUpdateFlag
	prevState   kvs.ValueState
	prevOp      kvs.TxnOperation
	prevErr     string
	prevDetails []string

	// for ExecuteTxnOperation
	dryExec   bool      // in
	execStart time.Time // out
	execStop  time.Time // out
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
			Source: exec.KVSource{
				Origin:      kv.origin,
				DerivedFrom: "", // base value
			},
			Context: &kvContext{
				isRetry: txn.txnType == kvs.RetryFailedOps,
			},
		})
	}

	s.execEngine.RunTransaction(txnCtx, kvChanges)

	// get rid of uninteresting intermediate pending Create/Delete operations
	executed = s.compressTxnOps(txnCtx.executed)
	return executed
}

// PrepareTxnOperation updates the underlying graph to reflect the requested
// value change.
// It also possible to skip operation execution and order the execution engine
// to move directly to Finalization without interruption.
// <prevValue> is the previous value as set by the given source (i.e. for merged
// value it is the previous value only for that single source) - ignored if
// <KeepValue> is enabled.
func (s *Scheduler) PrepareTxnOperation(txnPrivCtx exec.OpaqueCtx, kv *exec.KVChange, isRevert bool) (
	skipExec bool, prevValue proto.Message) {

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
	sources := getNodeSources(kvCtx.node)
	if !kv.KeepValue {
		prevValue = sources.GetSourceValue(kv.Source)
	}

	// update value sources
	if !kv.KeepValue {
		if kv.NewValue == nil {
			sources = sources.WithoutSource(kv.Source, !txnCtx.dryRun)
		} else {
			sources = sources.WithSource(ValueSource{
				KVSource: kv.Source,
				Value:    kv.NewValue,
			}, !txnCtx.dryRun)
		}
		kvCtx.node.SetFlags(sources)

		// SB is overshadowed by NB
		if len(sources.GetSources()) > 0 && !sources.IsObtained() &&
			kv.Source.Origin == kvs.FromSB {
			skipExec = true
			return
		}
	}

	// remember previous value status to detect and notify about changes
	kvCtx.prevUpdate = getNodeLastUpdate(kvCtx.node)
	kvCtx.prevState = getNodeState(kvCtx.node)
	kvCtx.prevOp = getNodeLastOperation(kvCtx.node)
	kvCtx.prevErr = getNodeErrorString(kvCtx.node)
	kvCtx.prevDetails = getValueDetails(kvCtx.node)

	// prepare operation description - fill attributes that we can even before
	// executing the operation
	kvCtx.opRecord = s.preRecordTxnOp2(kv, kvCtx, isRevert)

	// merge all value sources into one value
	for _, source := range sources.GetSources() {
		if !sources.IsObtained() && source.Origin == kvs.FromSB {
			// SB is overshadowed by NB
			continue
		}
		value := source.Value
		if value != nil {
			if kvCtx.newMergedValue == nil {
				kvCtx.newMergedValue = proto.Clone(value)
			} else {
				proto.Merge(kvCtx.newMergedValue, value)
			}
		}
	}

	// determine the operation type
	if kv.KeepValue {
		kvCtx.operation = kvs.TxnOperation_UNDEFINED // determined in IsTxnOperationReady
	} else if kvCtx.newMergedValue == nil {
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
		lastUpdateFlag.value = kvCtx.newMergedValue
	}
	if txnCtx.txnType == kvs.NBTransaction {
		lastUpdateFlag.retryEnabled = txnCtx.nb.retryEnabled
		lastUpdateFlag.retryArgs = txnCtx.nb.retryArgs
	} else if kvCtx.prevUpdate != nil {
		// inherit retry arguments from the last NB txn for this value
		lastUpdateFlag.retryEnabled = kvCtx.prevUpdate.retryEnabled
		lastUpdateFlag.retryArgs = kvCtx.prevUpdate.retryArgs
	} else if !kv.KeepValue && kv.Source.DerivedFrom != "" {
		// inherit from the parent value
		parentNode := txnCtx.graphW.GetNode(kv.Source.DerivedFrom)
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
	node.SetValue(kvCtx.newMergedValue)

	// add descriptor flag
	if !kvCtx.descriptor.isNil() {
		node.SetFlags(&DescriptorFlag{kvCtx.descriptor.name()})
		node.SetLabel(kvCtx.descriptor.keyLabel(kv.Key))
	}

	// handle unimplemented value
	sources := getNodeSources(kvCtx.node)
	unimplemented := !sources.IsObtained() && !sources.IsDerivedOnly() && kvCtx.descriptor.isNil()
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
		s.updateNodeState(node, kvCtx.opRecord.NewState, nil)
		return
	}

	// validate value
	if !txnCtx.dryRun && !sources.IsObtained() {
		err := kvCtx.descriptor.validate(node.GetKey(), node.GetValue())
		if err != nil {
			skipExec = true
			node.SetFlags(&UnavailValueFlag{})
			kvCtx.opRecord.NewErr = err
			kvCtx.opRecord.NewState = kvs.ValueState_INVALID
			kvCtx.opRecord.NOOP = true
			s.updateNodeState(node, kvCtx.opRecord.NewState, nil)
			node.SetFlags(&ErrorFlag{err: err, retriable: false})
			return
		}
	}

	// apply new relations
	derivedVals := kvCtx.descriptor.derivedValues(node.GetKey(), node.GetValue())
	dependencies := kvCtx.descriptor.dependencies(node.GetKey(), node.GetValue())
	node.SetTargets(constructTargets(dependencies, derivedVals))

	if sources.IsObtained() {
		// nothing to execute for SB notifications
		skipExec = true
		return
	}

	// continue with IsTxnOperationReady...
	return
}

func (s *Scheduler) prepareForUpdate(txnCtx *txnContext, kv *exec.KVChange, kvCtx *kvContext, isRevert bool) {

	// TODO
	return
}

// IsTxnOperationReady should determine whether to:
//  - proceed with operation execution or skip the execution and order
//    the execution engine to move to Finalization without any (additional)
//    interruption.
//  - wait for some other key-value pairs (of the same transaction) to
//    be changed/checked first - once those values are finalized, the readiness
//    check is repeated (skipping the values from precededBy which were already
//    processed) and the value change process continues accordingly
//  - block (freeze) some other values from entering the state machine while
//    this value is waiting/being executed (unfrozen when finalized) - if values
//    to be frozen are still in-progress, they will be finalized first (unless
//    this operation is preceding them)
func (s *Scheduler) IsTxnOperationReady(txnPrivCtx exec.OpaqueCtx, kv *exec.KVChange) (
	skipExec bool, precededBy []exec.KVChange, freeze utils.KeySet) {

	//txnCtx := txnPrivCtx.(*txnContext)
	kvCtx := kv.Context.(*kvContext)
	node := kvCtx.node

	// determine operation for a dependency update
	if kvCtx.operation == kvs.TxnOperation_UNDEFINED {
		kvCtx.operation = s.determineDepUpdateOperation(kvCtx.node)
		getNodeLastUpdate(kvCtx.node).txnOp = kvCtx.operation
		if kvCtx.operation == kvs.TxnOperation_UNDEFINED {
			// nothing to update
			skipExec = true
			return
		}
	}

	switch kvCtx.operation {
	case kvs.TxnOperation_DELETE:
		// TODO

	case kvs.TxnOperation_CREATE:
		// check if dependencies are available
		// Notes:
		//   - nodes with in-progress Create operations also appear available
		//   - dependencies will get frozen and once in-progress Create operation
		//     finalize, IsTxnOperationReady will get re-run
		if !isNodeReady(node) {
			// if not ready, nothing to do
			skipExec = true
			node.SetFlags(&UnavailValueFlag{})
			node.DelFlags(ErrorFlagIndex)
			kvCtx.opRecord.NewState = kvs.ValueState_PENDING
			kvCtx.opRecord.NOOP = true
			s.updateNodeState(node, kvCtx.opRecord.NewState, nil)
			return
		}

		// freeze dependencies
		freeze := utils.NewSliceBasedKeySet()
		for _, depPerLabel := range node.GetTargets(DependencyRelation) {
			for _, depNode := range depPerLabel.Nodes {
				freeze.Add(depNode.GetKey())
			}
		}

	case kvs.TxnOperation_UPDATE:
		// TODO
	}


	return false, nil, nil
}

// ExecuteTxnOperation is run from another go routine by one of the workers.
// The method should apply the value change (call Create/Delete/Update on the
// associated descriptor) and return error if the operation failed.
// Special error value AsyncExecError can be returned to signal the execution
// engine that the operation will continue in the background (e.g. due
// to a blocking action) and should be resumed (repeated with the checkpoint
// given in the error) once signaled through the ResumeAsyncOperation()
// method of the execution engine.
func (s *Scheduler) ExecuteTxnOperation(workerID int, checkpoint int, kv *exec.KVChange) (err error) {
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
// (followedBy).
// It is also possible to request the engine to trigger the transaction revert.
func (s *Scheduler) FinalizeTxnOperation(txnPrivCtx exec.OpaqueCtx, kv *exec.KVChange,
	wasRevert bool, opRetval error) (revertTxn bool, followedBy []exec.KVChange) {

	txnCtx := txnPrivCtx.(*txnContext)
	kvCtx := kv.Context.(*kvContext)
	node := kvCtx.node
	sources := getNodeSources(node)

	if kvCtx.operation == kvs.TxnOperation_UNDEFINED {
		// nothing has been done
		return
	}

	// update metadata
	if kvCtx.operation != kvs.TxnOperation_DELETE {
		if !sources.IsDerivedOnly() && kvCtx.descriptor.withMetadata() {
			node.SetMetadataMap(kvCtx.descriptor.name())
			node.SetMetadata(kvCtx.newMetadata)
		}
	}

	// finalize operation
	switch kvCtx.operation {
	case kvs.TxnOperation_DELETE:
		// TODO

	case kvs.TxnOperation_CREATE:
		if !kvCtx.opRecord.NOOP {
			if opRetval == nil {
				// value successfully created
				node.DelFlags(ErrorFlagIndex, UnavailValueFlagIndex)
				if sources.IsObtained() {
					kvCtx.opRecord.NewState = kvs.ValueState_OBTAINED
				} else {
					kvCtx.opRecord.NewState = kvs.ValueState_CONFIGURED
				}
				s.updateNodeState(node, kvCtx.opRecord.NewState, nil)
				followedBy = append(followedBy, s.scheduleDepUpdates(kvCtx, true)...)
				followedBy = append(followedBy, s.scheduleDerivedUpdates(kvCtx, false)...)
			} else {
				// execution ended with error
				node.SetFlags(&UnavailValueFlag{})
				retriableErr := kvCtx.descriptor.isRetriableFailure(opRetval)
				kvCtx.opRecord.NewErr = opRetval
				kvCtx.opRecord.NewState = s.markFailedValue2(
					node, opRetval, wasRevert, txnCtx, retriableErr)
			}
		}

	case kvs.TxnOperation_UPDATE:
		// TODO
	}

	// detect value state changes
	if !txnCtx.dryRun {
		nodeR := txnCtx.graphW.GetNode(kv.Key)
		stateChanged := kvCtx.prevUpdate == nil
		stateChanged = stateChanged || kvCtx.prevState != getNodeState(nodeR)
		stateChanged = stateChanged || kvCtx.prevOp != getNodeLastOperation(nodeR)
		stateChanged = stateChanged || kvCtx.prevErr != getNodeErrorString(nodeR)
		stateChanged = stateChanged || !equalValueDetails(kvCtx.prevDetails, getValueDetails(nodeR))
		if stateChanged {
			// update status of all base values from which this value is derived from
			for _, baseNode := range getNodeBaseSources(nodeR, txnCtx.graphW) {
				s.updatedStates.Add(baseNode.GetKey())
			}
		}
	}

	// TODO: do not revert on invalid value not originating from this transaction
	withRevert := txnCtx.txnType == kvs.NBTransaction && txnCtx.nb.revertOnFailure
	revertTxn = withRevert && kvCtx.opRecord.NewErr != nil

	// record the operation
	txnCtx.executed = append(txnCtx.executed, kvCtx.opRecord)
	return
}

// scheduleDepUpdates prepares a list of key-value pairs which will be re-check
// for dependencies.
func (s *Scheduler) scheduleDepUpdates(kvCtx *kvContext, forUnavailable bool) (updates []exec.KVChange) {
	var depNodes []graph.Node
	node := kvCtx.node
	for _, depPerLabel := range node.GetSources(DependencyRelation) {
		depNodes = append(depNodes, depPerLabel.Nodes...)
	}

	// order depNodes by key (just for deterministic behaviour which simplifies testing)
	sort.Slice(depNodes, func(i, j int) bool { return depNodes[i].GetKey() < depNodes[j].GetKey() })

	for _, depNode := range depNodes {
		if getNodeSources(depNode).IsObtained() {
			continue
		}
		if !isNodeAvailable(depNode) != forUnavailable {
			continue
		}

		updates = append(updates, exec.KVChange{
			Key:       depNode.GetKey(),
			KeepValue: true,
			Context:   &kvContext{},
		})
	}
	return
}

// scheduleDerivedUpdates prepares a list of derived key-value pairs for value-change.
func (s *Scheduler) scheduleDerivedUpdates(kvCtx *kvContext, remove bool) (updates []exec.KVChange) {
	node := kvCtx.node
	derivedVals := kvCtx.descriptor.derivedValues(node.GetKey(), node.GetValue())

	// order derivedVals by key (just for deterministic behaviour which simplifies testing)
	sort.Slice(derivedVals, func(i, j int) bool { return derivedVals[i].Key < derivedVals[j].Key })

	for _, derived := range derivedVals {
		if derived.Value == nil {
			s.Log.WithFields(logging.Fields{
				"key":          derived.Key,
				"derived-from": node.GetKey(),
			}).Warn("Derived nil value")
			continue
		}

		value := derived.Value
		if remove {
			value = nil
		}
		origin := kvs.FromNB
		if getNodeSources(node).IsObtained() {
			origin = kvs.FromSB
		}
		updates = append(updates, exec.KVChange{
			Key:       derived.Key,
			NewValue:  value,
			Source:    exec.KVSource{
				Origin:      origin,
				DerivedFrom: node.GetKey(),
			},
			Context:   &kvContext{},
		})
	}
	return
}

// PrepareForTxnRevert is run before reverting of already applied key-value
// changes is started (due to error(s)).
func (s *Scheduler) PrepareForTxnRevert(txnPrivCtx exec.OpaqueCtx, failedKVChanges utils.KeySet) {
	txnCtx := txnPrivCtx.(*txnContext)
	s.refreshGraph(txnCtx.graphW, failedKVChanges, nil, true)
}

// markFailedValue2 (will replace markFailedValue) decides whether to retry failed
// operation or not and updates the node state accordingly.
func (s *Scheduler) markFailedValue2(node graph.NodeRW, err error, wasRevert bool,
	txnCtx *txnContext, retriableErr bool) (newState kvs.ValueState) {

	// decide value state between FAILED and RETRYING
	newState = kvs.ValueState_FAILED
	toBeReverted := txnCtx.txnType == kvs.NBTransaction && txnCtx.nb.revertOnFailure && !wasRevert
	if retriableErr && !toBeReverted {
		// consider operation retry
		var alreadyRetried bool
		if txnCtx.txnType == kvs.RetryFailedOps {
			// TODO: handle multi-source
			baseKey := getNodeBaseKey(node)
			_, alreadyRetried = txnCtx.retry.keys[baseKey]
		}
		attempt := 1
		if alreadyRetried {
			attempt = txnCtx.retry.attempt + 1
		}
		lastUpdate := getNodeLastUpdate(node)
		if lastUpdate.retryEnabled && lastUpdate.retryArgs != nil &&
			(lastUpdate.retryArgs.MaxCount == 0 || attempt <= lastUpdate.retryArgs.MaxCount) {
			// retry is allowed
			newState = kvs.ValueState_RETRYING
		}
	}
	s.updateNodeState(node, newState, nil)
	node.SetFlags(&ErrorFlag{err: err, retriable: retriableErr})
	return newState
}
