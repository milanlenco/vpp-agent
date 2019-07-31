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
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/ligato/vpp-agent/plugins/kvscheduler/internal/utils"
	kvs "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
)

// KVChange represents a key-value pair which is being / going to be changed
// by the transaction.
type KVChange struct {
	Key       string
	// Note: multiple sources are allowed to change (subset of) the value
	// at the same time. But typically this is a one-item array.
	NewValues []ValueWithSource

	// Keep the current value, just re-check the state (dependencies, etc.) and
	// act accordingly. If enabled, NewValues should be ignored by TxnExecHandler.
	KeepValue bool

	// Private context to be used by TxnExecHandler.
	Context OpaqueCtx // TODO: how to merge kvContext of multiple KVChanges for the same value (but different sources; not KeepValue)
}

// ValueSource determines where the given value came from.
// Possibilities are:
//   - value requested by NB (NB value)
//   - value derived from a NB value
//   - value received as notification from SB (SB notification)
//   - value derived from a SB notification
type ValueSource struct {
	Origin      kvs.ValueOrigin
	DerivedFrom string
}

// ValueWithSource associates value with its source.
// Note: value is allowed to have multiple sources. Values from different
// sources applied for the same key should be merged into one by the execution
// handler.
type ValueWithSource struct {
	ValueSource
	Value proto.Message
}

// OpaqueCtx is used to plug arbitrary data into the execution context.
type OpaqueCtx interface{}

// AsyncExecError can be returned by ExecuteTxnOperation to tell the engine
// that the executed operation will continue asynchronously in the background.
// The operation execution will be resumed (repeated with the given checkpoint)
// once ResumeAsyncOperation is called for the given key.
type AsyncExecError struct {
	Checkpoint int
}

// Error returns a string representation of AsyncExecError.
func (e *AsyncExecError) Error() string {
	return fmt.Sprintf("operation continues asynchronously after reaching checkpoint: %d",
		e.Checkpoint)
}


// TxnExecEngine is the interface of the transaction execution engine.
// It is a state machine on top of which the transaction operation scheduling is
// operated by KVScheduler.
// A new instance of the engine is created with NewTxnExecEngine(). The constructor
// immediately starts the requested number of worker go routines. The workers
// remain idle, waiting for transaction operations to execute, and get stopped
// when the engine is closed with Close().
// To execute a new transaction, call RunTransaction. The method is synchronous
// (i.e. blocking), but the requested key-value changes are load-balanced across
// workers and potentially executed asynchronously to maximize the utilization.
//
// The engine is used to help with the graph walk without even understanding the
// details of the graph representation/implementation at this abstraction level.
// Change of every key-value pair (graph node visit) is performed in 4 steps:
//  1. Preparation
//      - the underlying handler is supposed to update the corresponding
//        graph node and the attached edges
//      - opaque context attached to key-value can be used to propagate further
//        input arguments for the execution into the worker (e.g. operation to
//        execute)
//  2. Readiness-check:
//      - determine whether to:
//          - proceed with operation execution
//          - skip operation execution (skip directly to step 4.)
//          - wait for some other key-value pairs (of the same transaction) to
//            be changed first - once those values are finalized, the readiness
//            check is repeated and the value change process continues accordingly
//          - block (freeze) some other values from entering the state machine while
//            this value is waiting/being executed (unfrozen when finalized) - if
//            values to be frozen are still in-progress, they will be finalized
//            first (unless this operation is preceding them)
//  3. Execution:
//      - run by one of the workers (i.e. different go routine)
//      - the underlying handler is supposed to execute the given operation
//      - the operation execution may be even send into the background to continue
//        asynchronously (typically to avoid the worker go routine to be suspended
//        while waiting for an external event) and resumed to continue from the
//        given checkpoint
//  4. Finalization:
//      - the underlying handler is supposed to:
//          - update the corresponding graph node to reflect the operation return
//            value
//          - determine if some more key-value pairs need to change as a consequence
//          - opaque context attached to the transaction can be used for example
//            to add recording of the operation, etc.
//          - the instance of the processed KVChange and the attached opaque
//            context are thrown away after the operation
//
// In principle, the steps 1-4 represent a graph node visit. The step 2 may cause
// the visit to be delayed until other nodes have been finalized, step 3 performs
// the actual operation and step 4 may enqueue adjacent nodes to be visited later.
// The next set of nodes to visit is added into the back of the queue, thus
// the graph is walked in the BFS order.
type TxnExecEngine interface {
	// RunTransaction executes transaction containing a given set of key-value
	// change requests.
	RunTransaction(txnCtx OpaqueCtx, input []KVChange)

	// ResumeAsyncOperation signals the execution engine to continue with
	// the execution for the ongoing asynchronous operation associated with the
	// given key.
	// if <done> is true, the operation will not be resumed and instead it will
	// be considered as done with the given error. In any case a non-nil <err>
	// stops the execution and the received error is forwarded to the main thread
	// of the execution handler to be processed in FinalizeTxnOperation.
	ResumeAsyncOperation(key string, done bool, err error)

	// Close stops all the worker go routines.
	Close() error
}

// TxnExecHandler implements the 4 steps of key-value change.
// PrepareTxnOperation, IsTxnOperationReady and FinalizeTxnOperation are run in
// the context of the main thread (thread from which RunTransaction was triggered),
// whereas calls to ExecuteTxnOperation are run in parallel, each assigned for
// execution to one of the workers.
type TxnExecHandler interface {
	//// KV Change Processing:

	// PrepareTxnOperation should update the underlying graph (abstracted-away
	// at the level of the execution engine) to reflect the requested value change.
	// Multiple values defined for the same key by different sources should be merged
	// into one proto message.
	// It is also possible to skip operation execution and order the execution engine
	// to move directly to Finalization without interruption.
	// <prevValues> should contain previous value for every source whose value
	// assigned to this key has been changed by the operation - ignored
	// if <KeepValue> is enabled. Since most values are typically single-source,
	// this is usually just a one-item slice.
	PrepareTxnOperation(txnCtx OpaqueCtx, kv *KVChange, isRevert bool) (
		skipExec bool, prevValues []ValueWithSource)

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
	IsTxnOperationReady(txnCtx OpaqueCtx, kv *KVChange) (
		skipExec bool, precededBy []KVChange, freeze utils.KeySet)

	// ExecuteTxnOperation is run from another go routine by one of the workers.
	// The method should apply the value change (call Create/Delete/Update on the
	// associated descriptor) and return error if the operation failed.
	// Special error value AsyncExecError can be returned to signal the execution
	// engine that the operation will continue in the background (e.g. due
	// to a blocking action) and should be resumed (repeated with the checkpoint
	// given in the error) once signaled through the ResumeAsyncOperation()
	// method of the execution engine.
	ExecuteTxnOperation(workerID, checkpoint int, kv *KVChange) (err error)

	// FinalizeTxnOperation is run after the operation has been executed/skipped.
	// Some more key-value pairs may be requested to be changed as a consequence
	// (followedBy).
	// It is also possible to request the engine to trigger the transaction revert.
	FinalizeTxnOperation(txnCtx OpaqueCtx, kv *KVChange, wasRevert bool, opRetval error) (
		revertTxn bool, followedBy []KVChange)

	//// Revert:

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
func (e *txnExecEngine) RunTransaction(txnCtx OpaqueCtx, input []KVChange) {
	// TODO
}

// ResumeAsyncOperation signals the execution engine to continue with
// the execution for the ongoing asynchronous operation associated with the
// given key.
// if <done> is true, the operation will not be resumed and instead it will
// be considered as done with the given error. In any case a non-nil <err>
// stops the execution and the received error is forwarded to the main thread
// of the execution handler to be processed in FinalizeTxnOperation.
func (e *txnExecEngine) ResumeAsyncOperation(key string, done bool, err error) {
	// TODO
}

// Close stops all the worker go routines.
func (e *txnExecEngine) Close() error {
	// TODO
	return nil
}


// TODO: couple of notes:
// - the value should not be in the queue more than once
//    - value change overwrites dependency update
//    - multiple planned dependency updates are pointless
//    - dependency update is basically already included in the value change
//    - multiple value changes (derived from multiple sources + combined with base)
//      can be merged into one operation
// - when transaction is started all value-change requests should be immediately
//   put into the queue - dependency-update before the first value-change will be
//   therefore omitted
// - when execution is skipped, FinalizeTxnOperation is called with nil opRetval
// - if the value is already being executed or is waiting, another operation
//   should be blocked - i.e. add "blocked" queue - and re-enqueued ASAP to the
//   front (regardless whether BFS or DFS is being used as it was already decided
//   to do that operation at the given moment, but it had to be blocked
//     - XXX not needed to behave like this, just enqueue as BFS or DFS and when
//       it is to-be processed, either postpone or go ahead)
// - whether to use DFS or BFS is from the functional point irrelevant, but BFS
//   will be more efficient since multiple derivations of the same value could be
//   merged into one operation
// - graph updates will be done in the preparation phase - that means, however,
//   that all the related values must be "blocked" (frozen - chose one of these words)
//    - those which are already being updated will be waited for (and not allowed to
//      re-enter), others will simply not be allowed to be enqueued
//    - frozen values will be put into the "blocked" queue
//    - frozen until the value that froze them is finalized
// - blocking will be used to cover:
//    - dependencies which are being executed
//    - dependencies which are waiting
//    - dependencies which are not in-progress, but could in theory overtake the values
//      that depend on them (slower worker)
//    - (not really needed) maybe freeze also parent values (i.e. implicit dependencies)
// - the relation is-preceded-by implies is-frozen-by, i.e:
//    - A -deps-> B -deps-> C
//    - delete of C is preceded by dep-update of B
//    - delete of B is preceded by dep-update of A + C is frozen (NOOP, already frozen implicitly)
//    - delete of A freezes B, C (both NOOP, already frozen implicitly - run transitive closure)
// - error during execution should be processed with priority
// - the size of the queue for the execution should be a multiple of the worker count
// - beware: ResumeAsyncOperation can overtake AsyncExecError (so on return the operation
//   could be immediately resumed)
// - it is possible to skipExec even from IsReady.. to handle cases when value
//   that was waited for has failed to get applied
// - Create/Update/Delete operation walk-through:
//    - Create
//       - Prepare:
//          - add node (without unavailable flag, pretending success until something fails)
//          - add relations
//       - IsReady:
//          - skip exec if some dependency is missing
//          - freeze dependencies
//       - Execute:
//          - Create the thing
//       - Finalize
//          - mark node as available
//          - followed up by the creation of derived values
//    - Delete
//       - Prepare:
//          - mark node as unavailable
//       - IsReady:
//          - preceded by removal of derived values and dependency check of values
//            that depend on it
//       - Execute:
//          - Delete the thing
//       - Finalize
//          - remove node from the graph
//    - Update (not re-create)
//       - Prepare:
//          - determine if equivalent and whether to re-create
//          - without changing relations, determine the set of obsolete, new
//            derived values and dependencies
//          - if some dependency is missing select Delete operation (write to context)
//             - mark node as unavailable
//       - IsReady:
//          - If to-be Updated
//             - freeze obsolete and new dependencies
//             - preceded by removal of post-update obsolete derived values
//          - else: (to-be deleted)
//             - preceded by removal of derived values and dependency check of values
//               that depend on it
//       - Execute:
//          - Delete or Update the thing
//       - Finalize
//          - update relations
//          - if was Update:
//             - followed up by the creation of derived values
//          - else (deleted)
//             - update relations
//    - Update with re-create
//       - 1st round:
//          - Prepare:
//              - determine if equivalent and whether to re-create
//              - mark node as unavailable
//          - IsReady:
//              - preceded by removal of derived values and dependency check of values
//                that depend on it
//          - Execute:
//              - Delete the thing
//          - Finalize
//              - followed by Create for this key and the new value
//       - 2nd round:
//          - Prepare:
//              - update relations
//              - skip exec if some dependency is missing
//          - IsReady:
//              - freeze dependencies
//          - Execute:
//              - Create the thing
//          - Finalize
//              - mark node as available
//              - followed up by the creation of derived values


// - TODO: unsolved issues
// - how to approach Refresh with multiple levels of derived values?
//   - e.g. for 2 interfaces in bridge domain, one of them not inserted into BD
//      - should Interconnection descriptor dump incomplete model?
//      - should Interconnection descriptor even implement Retrieve?
//      - how will refresh determine if the value actually exists (is dumped)
//        and is not just derived out
//        - I guess if the value has descriptor with Retrieve implemented, then
//          that Retrieve determines availability and actual value
//        - otherwise the value and availability is given by derivation, and SB/NB
//          base values without Retrieve remain unchanged
//
// - how to approach Retry?
//      - if merged base/derived-only/mix value fails, what to retry?
//      - should all base parents be retried?
//
// - how to approach Resync (efficiently)?
//      - merged value may be edited multiple times
//      - for example, consider multiple interconnect model instances deriving
//        (a subset) of the same bridge domain
//      - perhaps the kv-changes for the same value but from different sources
//        could be merged into one iteration of the execution engine
//
// - Observations:
//      - most likely the "last-update" will have to be per-source and included
//        in the value-sources flag
//      - Retrieve returns the merged value and overwrites whatever would be
//        a result of merging the derived values
//      - retry might need to be:
//          - per-source OR
//          - refresh+retry every closest parent value with Retrieve defined
//            on every branch following inverse of derivations (uh oh)
//      - Refresh would update:
//          - node actual value (Retrieve overwrites merge of derived)
//          - unavailability flag - NB values which are not only derived will be
//            marked as not available instead of deleted, retrieved/derived values
//            will be marked available
//          - value sources ? - probably YES to keep it in-sync with relations
//             - i.e. removes derives-from for obsolete derivations
//             - adds derived-from without last-update-txn info for derivations
//               which actually do exist
//      - BFS could be more efficient than DFS with high-level models and derived/merged
//        values if value changes are merged under one operation

// - Ideas:
//    - async goVPP:
//       func (ch *Channel) SendAsyncRequest(req, resp api.Message, receiveClb func(error)) {
//          ...
//       }
//   - async vppcalls:
//      func (h *InterfaceVppHandler) InterfaceAdminUp(ifIdx uint32, receiveClb func(error)) {
//          req := &interfaces.SwInterfaceSetFlags{
//              SwIfIndex:   ifIdx,
//              AdminUpDown: boolToUint(1),
//          }
//          reply := &interfaces.SwInterfaceSetFlagsReply{}
//          h.callsChannel.SendAsyncRequest(req, reply, receiveClb)
//      }
//



