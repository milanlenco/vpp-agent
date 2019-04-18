// Copyright (c) 2018 Cisco and/or its affiliates.
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
	"errors"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/gomega"

	. "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	"github.com/ligato/vpp-agent/plugins/kvscheduler/internal/test"
	"github.com/ligato/vpp-agent/plugins/kvscheduler/internal/utils"
)

func TestFailedDeleteOfDerivedValue(t *testing.T) {
	RegisterTestingT(t)

	// prepare KV Scheduler
	scheduler := NewPlugin(UseDeps(func(deps *Deps) {
		deps.HTTPHandlers = nil
	}))
	err := scheduler.Init()
	Expect(err).To(BeNil())

	// prepare mocks
	mockSB := test.NewMockSouthbound()
	// descriptor:
	descriptor := test.NewMockDescriptor(&KVDescriptor{
		Name:          descriptor1Name,
		NBKeyPrefix:   prefixA,
		KeySelector:   prefixSelector(prefixA),
		ValueTypeName: proto.MessageName(test.NewArrayValue()),
		DerivedValues: test.ArrayValueDerBuilder,
		WithMetadata:  true,
	}, mockSB, 0)
	scheduler.RegisterKVDescriptor(descriptor)

	// run non-resync transaction against empty SB
	schedulerTxn := scheduler.StartNBTransaction()
	schedulerTxn.SetValue(prefixA+baseValue1, test.NewArrayValue("item1"))
	seqNum, err := schedulerTxn.Commit(testCtx)
	Expect(seqNum).To(BeEquivalentTo(0))
	Expect(err).ShouldNot(HaveOccurred())

	// check the state of SB
	Expect(mockSB.GetKeysWithInvalidData()).To(BeEmpty())
	// -> base value 1
	value := mockSB.GetValue(prefixA + baseValue1)
	Expect(value).ToNot(BeNil())
	Expect(proto.Equal(value.Value, test.NewArrayValue("item1"))).To(BeTrue())
	Expect(value.Metadata).ToNot(BeNil())
	Expect(value.Metadata.(test.MetaWithInteger).GetInteger()).To(BeEquivalentTo(0))
	Expect(value.Origin).To(BeEquivalentTo(FromNB))
	// -> item1 derived from base value 1
	value = mockSB.GetValue(prefixA + baseValue1 + "/item1")
	Expect(value).ToNot(BeNil())
	Expect(proto.Equal(value.Value, test.NewStringValue("item1"))).To(BeTrue())
	Expect(value.Metadata).To(BeNil())
	Expect(value.Origin).To(BeEquivalentTo(FromNB))

	// plan error before 2nd txn
	failedDeleteClb := func() {
		mockSB.SetValue(prefixA+baseValue1, test.NewArrayValue("item1"),
			&test.OnlyInteger{Integer: 0}, FromNB, false)
	}
	mockSB.PlanError(prefixA+baseValue1+"/item1", errors.New("failed to delete value"), failedDeleteClb)

	// run 2nd non-resync transaction that will have errors
	startTime := time.Now()
	schedulerTxn2 := scheduler.StartNBTransaction()
	schedulerTxn2.SetValue(prefixA+baseValue1, nil)
	seqNum, err = schedulerTxn2.Commit(testCtx)
	stopTime := time.Now()
	Expect(seqNum).To(BeEquivalentTo(1))
	Expect(err).ToNot(BeNil())
	txnErr := err.(*TransactionError)
	Expect(txnErr.GetTxnInitError()).ShouldNot(HaveOccurred())
	kvErrors := txnErr.GetKVErrors()
	Expect(kvErrors).To(HaveLen(1))
	Expect(kvErrors[0].Key).To(BeEquivalentTo(prefixA + baseValue1 + "/item1"))
	Expect(kvErrors[0].TxnOperation).To(BeEquivalentTo(TxnOperation_DELETE))
	Expect(kvErrors[0].Error.Error()).To(BeEquivalentTo("failed to delete value"))

	// check transaction operations
	txnHistory := scheduler.GetTransactionHistory(time.Time{}, time.Now())
	Expect(txnHistory).To(HaveLen(2))
	txn := txnHistory[1]
	Expect(txn.PreRecord).To(BeFalse())
	Expect(txn.Start.After(startTime)).To(BeTrue())
	Expect(txn.Start.Before(txn.Stop)).To(BeTrue())
	Expect(txn.Stop.Before(stopTime)).To(BeTrue())
	Expect(txn.SeqNum).To(BeEquivalentTo(1))
	Expect(txn.TxnType).To(BeEquivalentTo(NBTransaction))
	Expect(txn.ResyncType).To(BeEquivalentTo(NotResync))
	Expect(txn.Description).To(BeEmpty())
	checkRecordedValues(txn.Values, []RecordedKVPair{
		{Key: prefixA + baseValue1, Value: nil, Origin: FromNB},
	})

	// -> planned
	txnOps := RecordedTxnOps{
		{
			Operation: TxnOperation_DELETE,
			Key:       prefixA + baseValue1 + "/item1",
			IsDerived: true,
			PrevValue: utils.RecordProtoMessage(test.NewStringValue("item1")),
			PrevState: ValueState_CONFIGURED,
			NewState:  ValueState_REMOVED,
		},
		{
			Operation: TxnOperation_DELETE,
			Key:       prefixA + baseValue1,
			PrevValue: utils.RecordProtoMessage(test.NewArrayValue("item1")),
			PrevState: ValueState_CONFIGURED,
			NewState:  ValueState_REMOVED,
		},
	}
	checkTxnOperations(txn.Planned, txnOps)

	// -> executed
	txnOps = RecordedTxnOps{
		{
			Operation: TxnOperation_DELETE,
			Key:       prefixA + baseValue1 + "/item1",
			IsDerived: true,
			PrevValue: utils.RecordProtoMessage(test.NewStringValue("item1")),
			PrevState: ValueState_CONFIGURED,
			NewState:  ValueState_FAILED,
			NewErr:    errors.New("failed to delete value"),
		},
	}
	checkTxnOperations(txn.Executed, txnOps)

	// check value status
	status := scheduler.GetValueStatus(prefixA + baseValue1)
	Expect(status).ToNot(BeNil())
	checkBaseValueStatus(status, &BaseValueStatus{
		Value: &ValueStatus{
			Key:           prefixA + baseValue1,
			State:         ValueState_CONFIGURED,
			LastOperation: TxnOperation_DELETE,
		},
		DerivedValues: []*ValueStatus{
			{
				Key:           prefixA + baseValue1 + "/item1",
				State:         ValueState_FAILED,
				LastOperation: TxnOperation_DELETE,
				Error:         "failed to delete value",
			},
		},
	})

	// close scheduler
	err = scheduler.Close()
	Expect(err).To(BeNil())
}

func TestFailedDeleteOfDependentValue(t *testing.T) {
	RegisterTestingT(t)

	// prepare KV Scheduler
	scheduler := NewPlugin(UseDeps(func(deps *Deps) {
		deps.HTTPHandlers = nil
	}))
	err := scheduler.Init()
	Expect(err).To(BeNil())

	// prepare mocks
	mockSB := test.NewMockSouthbound()
	// descriptors:
	descriptor1 := test.NewMockDescriptor(&KVDescriptor{
		Name:            descriptor1Name,
		NBKeyPrefix:     prefixA,
		KeySelector:     prefixSelector(prefixA),
		ValueTypeName:   proto.MessageName(test.NewStringValue("")),
		ValueComparator: test.StringValueComparator,
		UpdateWithRecreate: func(key string, oldValue, newValue proto.Message, metadata Metadata) bool {
			return true
		},
		Dependencies: func(key string, value proto.Message) []Dependency {
			if key == prefixA+baseValue1 {
				return []Dependency{
					{
						Label: prefixB,
						AnyOf: AnyOfDependency{
							KeyPrefixes: []string{prefixB},
						},
					},
				}
			}
			return nil
		},
		WithMetadata: true,
	}, mockSB, 0)
	descriptor2 := test.NewMockDescriptor(&KVDescriptor{
		Name:            descriptor2Name,
		NBKeyPrefix:     prefixB,
		KeySelector:     prefixSelector(prefixB),
		ValueTypeName:   proto.MessageName(test.NewStringValue("")),
		ValueComparator: test.StringValueComparator,
		UpdateWithRecreate: func(key string, oldValue, newValue proto.Message, metadata Metadata) bool {
			return true
		},
		WithMetadata: true,
	}, mockSB, 0)
	scheduler.RegisterKVDescriptor(descriptor1)
	scheduler.RegisterKVDescriptor(descriptor2)

	// get metadata map created for descriptor1
	metadataMap := scheduler.GetMetadataMap(descriptor1.Name)
	nameToInteger1, withMetadataMap := metadataMap.(test.NameToInteger)
	Expect(withMetadataMap).To(BeTrue())

	// run non-resync transaction against empty SB
	schedulerTxn := scheduler.StartNBTransaction()
	schedulerTxn.SetValue(prefixA+baseValue1, test.NewStringValue("A-item1"))
	schedulerTxn.SetValue(prefixB+baseValue1, test.NewStringValue("B-item1"))
	seqNum, err := schedulerTxn.Commit(testCtx)
	Expect(seqNum).To(BeEquivalentTo(0))
	Expect(err).ShouldNot(HaveOccurred())

	// check metadata
	metadata, exists := nameToInteger1.LookupByName(baseValue1)
	Expect(exists).To(BeTrue())
	Expect(metadata.GetInteger()).To(BeEquivalentTo(0))

	// plan error before 2nd txn
	// - delete will return error even though the value will be actually removed from SB
	failedDeleteClb := func() {
		mockSB.SetValue(prefixA+baseValue1, nil, nil, FromNB, false)
	}
	// first Re-Create will succeed
	mockSB.PlanError(prefixA+baseValue1, nil, nil)
	mockSB.PlanError(prefixA+baseValue1, nil, nil)
	mockSB.PlanError(prefixA+baseValue1, errors.New("failed to delete value"), failedDeleteClb)

	// run 2nd non-resync transaction that will have errors
	schedulerTxn2 := scheduler.StartNBTransaction()
	schedulerTxn2.SetValue(prefixA+baseValue1, test.NewStringValue("A-item1-rev2"))
	schedulerTxn2.SetValue(prefixB+baseValue1, test.NewStringValue("B-item1-rev2"))
	ctx := testCtx
	ctx = WithRetry(ctx, time.Second, 3, false)
	seqNum, err = schedulerTxn2.Commit(ctx)
	Expect(seqNum).To(BeEquivalentTo(1))
	Expect(err).ToNot(BeNil())
	txnErr := err.(*TransactionError)
	Expect(txnErr.GetTxnInitError()).ShouldNot(HaveOccurred())
	kvErrors := txnErr.GetKVErrors()
	Expect(kvErrors).To(HaveLen(1))
	Expect(kvErrors[0].Key).To(BeEquivalentTo(prefixA + baseValue1))
	Expect(kvErrors[0].TxnOperation).To(BeEquivalentTo(TxnOperation_DELETE))
	Expect(kvErrors[0].Error.Error()).To(BeEquivalentTo("failed to delete value"))

	// check transaction operations
	txnHistory := scheduler.GetTransactionHistory(time.Time{}, time.Now())
	Expect(txnHistory).To(HaveLen(2))
	txn := txnHistory[1]
	Expect(txn.PreRecord).To(BeFalse())
	Expect(txn.SeqNum).To(BeEquivalentTo(1))
	Expect(txn.TxnType).To(BeEquivalentTo(NBTransaction))
	Expect(txn.ResyncType).To(BeEquivalentTo(NotResync))
	Expect(txn.Description).To(BeEmpty())
	checkRecordedValues(txn.Values, []RecordedKVPair{
		{Key: prefixA + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")), Origin: FromNB},
		{Key: prefixB + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("B-item1-rev2")), Origin: FromNB},
	})

	// -> planned
	txnOps := RecordedTxnOps{
		{
			Operation:  TxnOperation_DELETE,
			Key:        prefixA + baseValue1,
			PrevValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1")),
			IsRecreate: true,
			PrevState:  ValueState_CONFIGURED,
			NewState:   ValueState_REMOVED,
		},
		{
			Operation:  TxnOperation_CREATE,
			Key:        prefixA + baseValue1,
			NewValue:   utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			IsRecreate: true,
			PrevState:  ValueState_REMOVED,
			NewState:   ValueState_CONFIGURED,
		},
		{
			Operation: TxnOperation_DELETE,
			Key:       prefixA + baseValue1,
			PrevValue: utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			PrevState: ValueState_CONFIGURED,
			NewState:  ValueState_PENDING,
		},
		{
			Operation:  TxnOperation_DELETE,
			Key:        prefixB + baseValue1,
			IsRecreate: true,
			PrevValue:  utils.RecordProtoMessage(test.NewStringValue("B-item1")),
			PrevState:  ValueState_CONFIGURED,
			NewState:   ValueState_REMOVED,
		},
		{
			Operation:  TxnOperation_CREATE,
			Key:        prefixB + baseValue1,
			IsRecreate: true,
			NewValue:   utils.RecordProtoMessage(test.NewStringValue("B-item1-rev2")),
			PrevState:  ValueState_REMOVED,
			NewState:   ValueState_CONFIGURED,
		},
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixA + baseValue1,
			PrevValue: utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			PrevState: ValueState_PENDING,
			NewState:  ValueState_CONFIGURED,
		},
	}
	checkTxnOperations(txn.Planned, txnOps)

	// -> executed
	txnOps = RecordedTxnOps{
		{
			Operation:  TxnOperation_DELETE,
			Key:        prefixA + baseValue1,
			PrevValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1")),
			IsRecreate: true,
			PrevState:  ValueState_CONFIGURED,
			NewState:   ValueState_REMOVED,
		},
		{
			Operation:  TxnOperation_CREATE,
			Key:        prefixA + baseValue1,
			NewValue:   utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			IsRecreate: true,
			PrevState:  ValueState_REMOVED,
			NewState:   ValueState_CONFIGURED,
		},
		{
			Operation: TxnOperation_DELETE,
			Key:       prefixA + baseValue1,
			PrevValue: utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			PrevState: ValueState_CONFIGURED,
			NewState:  ValueState_RETRYING,
			NewErr:    errors.New("failed to delete value"),
		},
	}
	checkTxnOperations(txn.Executed, txnOps)

	// check metadata (should be unset by Refresh)
	_, exists = nameToInteger1.LookupByName(baseValue1)
	Expect(exists).To(BeFalse())

	// wait for retry
	time.Sleep(3 * time.Second)

	// check transaction operations
	txnHistory = scheduler.GetTransactionHistory(time.Time{}, time.Now())
	Expect(txnHistory).To(HaveLen(3))
	txn = txnHistory[2]
	Expect(txn.PreRecord).To(BeFalse())
	Expect(txn.SeqNum).To(BeEquivalentTo(2))
	Expect(txn.TxnType).To(BeEquivalentTo(RetryFailedOps))
	Expect(txn.ResyncType).To(BeEquivalentTo(NotResync))
	Expect(txn.Description).To(BeEmpty())
	checkRecordedValues(txn.Values, []RecordedKVPair{
		{Key: prefixA + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")), Origin: FromNB},
		{Key: prefixB + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("B-item1-rev2")), Origin: FromNB},
	})

	txnOps = RecordedTxnOps{
		{
			Operation:  TxnOperation_DELETE,
			Key:        prefixB + baseValue1,
			IsRecreate: true,
			PrevValue:  utils.RecordProtoMessage(test.NewStringValue("B-item1")),
			PrevState:  ValueState_CONFIGURED,
			NewState:   ValueState_REMOVED,
		},
		{
			Operation:  TxnOperation_CREATE,
			Key:        prefixB + baseValue1,
			IsRecreate: true,
			NewValue:   utils.RecordProtoMessage(test.NewStringValue("B-item1-rev2")),
			PrevState:  ValueState_REMOVED,
			NewState:   ValueState_CONFIGURED,
		},
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixA + baseValue1,
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1-rev2")),
			PrevState: ValueState_PENDING,
			NewState:  ValueState_CONFIGURED,
		},
	}
	checkTxnOperations(txn.Planned, txnOps)
	checkTxnOperations(txn.Executed, txnOps)

	// check metadata
	metadata, exists = nameToInteger1.LookupByName(baseValue1)
	Expect(exists).To(BeTrue())
	Expect(metadata.GetInteger()).To(BeEquivalentTo(2))

	// close scheduler
	err = scheduler.Close()
	Expect(err).To(BeNil())
}

func TestRetryValueTogetherWithPendingDependencies(t *testing.T) {
	RegisterTestingT(t)

	// prepare KV Scheduler
	scheduler := NewPlugin(UseDeps(func(deps *Deps) {
		deps.HTTPHandlers = nil
	}))
	err := scheduler.Init()
	Expect(err).To(BeNil())

	// prepare mocks
	mockSB := test.NewMockSouthbound()
	// descriptors:
	descriptor1 := test.NewMockDescriptor(&KVDescriptor{
		Name:            descriptor1Name,
		NBKeyPrefix:     prefixA,
		KeySelector:     prefixSelector(prefixA),
		ValueTypeName:   proto.MessageName(test.NewStringValue("")),
		ValueComparator: test.StringValueComparator,
		Dependencies: func(key string, value proto.Message) []Dependency {
			if key == prefixA+baseValue1 {
				return []Dependency{
					{
						Label: prefixB + baseValue1,
						Key:   prefixB + baseValue1,
					},
				}
			}
			return nil
		},
		WithMetadata: true,
	}, mockSB, 0)
	descriptor2 := test.NewMockDescriptor(&KVDescriptor{
		Name:            descriptor2Name,
		NBKeyPrefix:     prefixB,
		KeySelector:     prefixSelector(prefixB),
		ValueTypeName:   proto.MessageName(test.NewStringValue("")),
		ValueComparator: test.StringValueComparator,
		WithMetadata:    true,
	}, mockSB, 0)
	scheduler.RegisterKVDescriptor(descriptor1)
	scheduler.RegisterKVDescriptor(descriptor2)

	// plan error before txn
	// - create fails, but actually the value was applied correctly
	mockSB.PlanError(prefixB+baseValue1, errors.New("failed to add value"),
		func() {
			mockSB.SetValue(prefixB+baseValue1, test.NewStringValue("B-item1"),
				&test.OnlyInteger{Integer: 0}, FromNB, false)
		})

	// run resync transaction against empty SB
	ctx := WithResync(testCtx, FullResync, true)
	ctx = WithRetry(ctx, time.Second, 3, false)
	schedulerTxn := scheduler.StartNBTransaction()
	schedulerTxn.SetValue(prefixA+baseValue1, test.NewStringValue("A-item1"))
	schedulerTxn.SetValue(prefixB+baseValue1, test.NewStringValue("B-item1"))
	seqNum, err := schedulerTxn.Commit(ctx)
	Expect(seqNum).To(BeEquivalentTo(0))
	Expect(err).ToNot(BeNil())
	txnErr := err.(*TransactionError)
	Expect(txnErr.GetTxnInitError()).ShouldNot(HaveOccurred())
	kvErrors := txnErr.GetKVErrors()
	Expect(kvErrors).To(HaveLen(1))
	Expect(kvErrors[0].Key).To(BeEquivalentTo(prefixB + baseValue1))
	Expect(kvErrors[0].TxnOperation).To(BeEquivalentTo(TxnOperation_CREATE))
	Expect(kvErrors[0].Error.Error()).To(BeEquivalentTo("failed to add value"))

	// check transaction operations
	txnHistory := scheduler.GetTransactionHistory(time.Time{}, time.Now())
	Expect(txnHistory).To(HaveLen(1))
	txn := txnHistory[0]
	Expect(txn.PreRecord).To(BeFalse())
	Expect(txn.SeqNum).To(BeEquivalentTo(0))
	Expect(txn.TxnType).To(BeEquivalentTo(NBTransaction))
	Expect(txn.ResyncType).To(BeEquivalentTo(FullResync))
	Expect(txn.Description).To(BeEmpty())
	checkRecordedValues(txn.Values, []RecordedKVPair{
		{Key: prefixA + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("A-item1")), Origin: FromNB},
		{Key: prefixB + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("B-item1")), Origin: FromNB},
	})

	// -> planned
	txnOps := RecordedTxnOps{
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixB + baseValue1,
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("B-item1")),
			PrevState: ValueState_NONEXISTENT,
			NewState:  ValueState_CONFIGURED,
		},
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixA + baseValue1,
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1")),
			PrevState: ValueState_NONEXISTENT,
			NewState:  ValueState_CONFIGURED,
		},
	}
	checkTxnOperations(txn.Planned, txnOps)

	// -> executed
	txnOps = RecordedTxnOps{
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixA + baseValue1,
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1")),
			PrevState: ValueState_NONEXISTENT,
			NewState:  ValueState_PENDING,
			NOOP:      true,
		},
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixB + baseValue1,
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("B-item1")),
			PrevState: ValueState_NONEXISTENT,
			NewState:  ValueState_RETRYING,
			NewErr:    errors.New("failed to add value"),
		},
	}
	checkTxnOperations(txn.Executed, txnOps)

	// wait for retry
	time.Sleep(3 * time.Second)

	// check transaction operations
	txnHistory = scheduler.GetTransactionHistory(time.Time{}, time.Now())
	Expect(txnHistory).To(HaveLen(2))
	txn = txnHistory[1]
	Expect(txn.PreRecord).To(BeFalse())
	Expect(txn.SeqNum).To(BeEquivalentTo(1))
	Expect(txn.TxnType).To(BeEquivalentTo(RetryFailedOps))
	Expect(txn.ResyncType).To(BeEquivalentTo(NotResync))
	Expect(txn.Description).To(BeEmpty())
	checkRecordedValues(txn.Values, []RecordedKVPair{
		{Key: prefixA + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("A-item1")), Origin: FromNB},
		{Key: prefixB + baseValue1, Value: utils.RecordProtoMessage(test.NewStringValue("B-item1")), Origin: FromNB},
	})

	txnOps = RecordedTxnOps{
		{
			Operation: TxnOperation_CREATE,
			Key:       prefixA + baseValue1,
			NewValue:  utils.RecordProtoMessage(test.NewStringValue("A-item1")),
			PrevState: ValueState_NONEXISTENT,
			NewState:  ValueState_CONFIGURED,
		},
	}
	checkTxnOperations(txn.Planned, txnOps)
	checkTxnOperations(txn.Executed, txnOps)

	// close scheduler
	err = scheduler.Close()
	Expect(err).To(BeNil())
}
