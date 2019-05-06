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

package exec_engine

import (
	"github.com/gogo/protobuf/proto"
	"github.com/ligato/vpp-agent/plugins/kvscheduler/internal/utils"
)

// KVChange represents a key-value pair which is being / going to be changed
// by the transaction.
type KVChange struct {
	Key       string
	NewValue  proto.Message
	Context   OpaqueCtx
}

// OpaqueCtx is used to plug arbitrary data into the execution context.
type OpaqueCtx interface{}

// TxnExecEngine is the interface of the transaction execution engine.
// A new instance of the engine is created with NewTxnExecEngine(). The constructor
// immediately starts the requested number of worker go routines. The workers
// remain idle, waiting for transaction operations to execute, and get stopped
// when the engine is closed with Close().
// To execute a new transaction, call RunTransaction. The method is synchronous
// (i.e. blocking), but the requested key-value changes are load-balanced across
// workers to maximize the utilization.
//
// The engine is used to help with the graph walk without even understanding the
// details of the graph representation/implementation at this abstraction level.
// Change of every key-value (graph node visit) pair is performed in 3 steps:
//  1. Preparation: the underlying handler is supposed to
//      - update the corresponding graph node and the attached edges
//     (Continuation:)
//      - determine whether to:
//          - proceed with operation execution
//          - skip operation execution (skip directly to step 3.)
//          - wait for some other key-value pairs (of the same transaction) to
//            be changed first - once those values are finalized, the value
//            change continues with ContinueTxnOperation
//      - opaque context attached to key-value can be used to propagate further
//        input arguments for the execution into the worker (e.g. operation to
//        execute)
// 2. Execution: run by one of the workers (i.e. different go routine)
//      - the underlying handler is supposed to execute the given operation
// 3. Finalization: the underlying handler is supposed to
//      - update the corresponding graph node to reflect the operation return value
//      - determine if some more key-value pair need to change as a consequence
//      - opaque context attached to transaction can be used for example to add
//        recording of the operation, etc.
//      - the instance of KVChange and the attached opaque context are thrown
//        away after the operation
//
// In principle, the steps 1-3 represent a graph node visit. The step 1 may cause
// the visit to be delayed until other nodes have been finalized, step 2 performs
// the actual operation and step 3 may enqueue adjacent nodes to be visited later.
// The next set of nodes to visit is added into the back of the queue, thus
// the graph is walked in the BFS order.
type TxnExecEngine interface {
	// RunTransaction executes transaction containing a given set of key-value
	// change requests.
	RunTransaction(txnCtx OpaqueCtx, input []KVChange, withRevert bool)

	// Close stops all the worker go routines.
	Close() error
}

// TxnExecHandler implements the 3 steps of key-value change.
// PrepareTxnOperation and FinalizeTxnOperation are run in the context of the
// main thread (thread from which RunTransaction was triggered), whereas calls
// to ExecuteTxnOperation are run in parallel, each assigned for execution to one
// of the workers.
type TxnExecHandler interface {
	// PrepareTxnOperation should determine whether to:
	//  - proceed with operation execution
	//  - skip operation execution (skip directly to Finalization)
	//  - wait for some other key-value pairs (of the same transaction) to
	//    be changed first (when those values are finalized, preparation for this
	//    value change will be replayed) -- TODO: avoid full replay
	// Furthermore, the underlying graph (abstracted-away at this level) should
	// be updated to reflect the value change.
	PrepareTxnOperation(txnCtx OpaqueCtx, kv *KVChange, isRevert bool) (
		prevValue proto.Message, skipExec bool,	waitFor utils.KeySet, err error)

	// TODO: maybe 4 steps after all?
	ContinueTxnOperation(txnCtx OpaqueCtx, kv *KVChange, isRevert bool, depOpRetval map[string]error) (
		skipExec bool, waitFor utils.KeySet, err error)

	// ExecuteTxnOperation is run from another go routine by one of the workers.
	// The method should apply the value change (call Create/Delete/Update on the
	// associated descriptor) and return error if the operation failed.
	ExecuteTxnOperation(workerID int, kv *KVChange) (err error)

	// FinalizeTxnOperation is run after the operation has been executed/skipped.
	// Some more key-value pair may be requested to be changed as a consequence
	// (follow-ups).
	FinalizeTxnOperation(txnCtx OpaqueCtx, kv *KVChange, wasRevert bool, opRetval error) (
		followUp []KVChange)

	// PrepareForTxnRevert is run before reverting of already applied key-value
	// changes is started (due to error(s)).
	PrepareForTxnRevert(txnCtx OpaqueCtx, failedKVChanges utils.KeySet)
}

// txnExecEngine implements the transaction execution handler.
type txnExecEngine struct{
	handler      TxnExecHandler
	numOfWorkers int
	verboseLog   bool

	// TODO
}

// NewTxnExecEngine is a constructor for the transaction execution engine.
func NewTxnExecEngine(handler TxnExecHandler, numOfWorkers int, verboseLog bool) TxnExecEngine {
	return &txnExecEngine{
		handler:      handler,
		numOfWorkers: numOfWorkers,
		verboseLog:   verboseLog,
	}
}

// RunTransaction executes transaction containing a given set of key-value
// change requests.
func (e *txnExecEngine) RunTransaction(txnCtx OpaqueCtx, input []KVChange, withRevert bool) {
	// TODO
}

// Close stops all the worker go routines.
func (e *txnExecEngine) Close() error {
	// TODO
	return nil
}