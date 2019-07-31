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
	"fmt"

	"github.com/gogo/protobuf/proto"

	kvs "github.com/ligato/vpp-agent/plugins/kvscheduler/api"
	exec "github.com/ligato/vpp-agent/plugins/kvscheduler/internal/exec-engine"
)

const (
	////// updated by transactions:

	// LastUpdateFlagName is the name of the LastUpdate flag.
	LastUpdateFlagName = "last-update"
	// LastUpdateFlagIndex is the Index of the LastUpdate flag.
	LastUpdateFlagIndex = 0

	// ErrorFlagName is the name of the Error flag.
	ErrorFlagName = "error"
	// ErrorFlagIndex is the Index of the ErrorFlag flag.
	ErrorFlagIndex = 1

	////// updated by transactions + refresh:

	// ValueStateFlagName is the name of the Value-State flag.
	ValueStateFlagName = "value-state"
	// ValueStateFlagIndex is the index of the Value-State flag.
	ValueStateFlagIndex = 2

	// UnavailValueFlagName is the name of the Unavailable-Value flag.
	UnavailValueFlagName = "unavailable"
	// UnavailValueFlagIndex is the index of the Unavailable-Value flag.
	UnavailValueFlagIndex = 3

	// DescriptorFlagName is the name of the Descriptor flag.
	DescriptorFlagName = "descriptor"
	// DescriptorFlagIndex is the index of the Descriptor flag.
	DescriptorFlagIndex = 4

	// TODO: remove the derived flag after migrating to execution engine
	// DerivedFlagName is the name of the Derived flag.
	DerivedFlagName = "derived"
	// DerivedFlagIndex is the index of the Derived flag.
	DerivedFlagIndex = 5

	// ValueSourceFlagName is the name of the Value-Source flag.
	ValueSourceFlagName = "value-source"
	// ValueSourceFlagIndex is the index of the Value-Source flag.
	ValueSourceFlagIndex = 6
)

// flagNameToIndex converts flag name to the associated index.
func flagNameToIndex(flagName string) int {
	switch flagName {
	case LastUpdateFlagName:
		return LastUpdateFlagIndex
	case ErrorFlagName:
		return ErrorFlagIndex
	case ValueStateFlagName:
		return ValueStateFlagIndex
	case UnavailValueFlagName:
		return UnavailValueFlagIndex
	case DescriptorFlagName:
		return DescriptorFlagIndex
	case DerivedFlagName:
		return DerivedFlagIndex
	case ValueSourceFlagName:
		return ValueSourceFlagIndex
	}
	return -1
}

/****************************** LastUpdate Flag *******************************/

// LastUpdateFlag is set to remember the last transaction which has
// changed/updated the value.
// Not set to values just discovered by refresh (state = DISCOVERED).
type LastUpdateFlag struct {
	txnSeqNum uint64
	txnOp     kvs.TxnOperation
	value     proto.Message  // the last applied *merged* value  TODO: needed? (currently used for downstream resync + retry + verification

	// updated only when the value content is being modified
	revert bool // TODO: needed?

	// set by NB txn, inherited by Retry and SB notifications
	retryEnabled bool
	retryArgs    *kvs.RetryOpt
}

// GetIndex returns 0.
func (flag *LastUpdateFlag) GetIndex() int {
	return LastUpdateFlagIndex
}

// GetName return name of the LastUpdate flag.
func (flag *LastUpdateFlag) GetName() string {
	return LastUpdateFlagName
}

// GetValue describes the last update (txn-seq number only).
func (flag *LastUpdateFlag) GetValue() string {
	return fmt.Sprintf("TXN-%d", flag.txnSeqNum)
}

/******************************* Error Flag ***********************************/

// ErrorFlag is used to store error returned from the last operation, including
// validation errors.
type ErrorFlag struct {
	err       error
	retriable bool
}

// GetIndex returns 1.
func (flag *ErrorFlag) GetIndex() int {
	return ErrorFlagIndex
}

// GetName return name of the Origin flag.
func (flag *ErrorFlag) GetName() string {
	return ErrorFlagName
}

// GetValue returns the error as string.
func (flag *ErrorFlag) GetValue() string {
	if flag.err == nil {
		return ""
	}
	return flag.err.Error()
}

/***************************** Value State Flag *******************************/

// ValueStateFlag stores current state of the value.
// Assigned to every value.
type ValueStateFlag struct {
	valueState kvs.ValueState
}

// GetIndex returns 2.
func (flag *ValueStateFlag) GetIndex() int {
	return ValueStateFlagIndex
}

// GetName returns name of the ValueState flag.
func (flag *ValueStateFlag) GetName() string {
	return ValueStateFlagName
}

// GetValue returns the string representation of the state.
func (flag *ValueStateFlag) GetValue() string {
	return flag.valueState.String()
}

/************************** Unavailable Value Flag ****************************/

// UnavailValueFlag is used to mark NB values which should not be considered
// when resolving dependencies of other values (for various possible reasons).
type UnavailValueFlag struct {
}

// GetIndex returns 3.
func (flag *UnavailValueFlag) GetIndex() int {
	return UnavailValueFlagIndex
}

// GetName return name of the UnavailValue flag.
func (flag *UnavailValueFlag) GetName() string {
	return UnavailValueFlagName
}

// GetValue return empty string (presence of the flag is the only information).
func (flag *UnavailValueFlag) GetValue() string {
	return ""
}

/*************************** Descriptor Value Flag ****************************/

// DescriptorFlag is used to lookup values by their descriptor.
// Not assigned to properties and UNIMPLEMENTED values.
type DescriptorFlag struct {
	descriptorName string
}

// GetIndex returns 4.
func (flag *DescriptorFlag) GetIndex() int {
	return DescriptorFlagIndex
}

// GetName return name of the Descriptor flag.
func (flag *DescriptorFlag) GetName() string {
	return DescriptorFlagName
}

// GetValue returns the descriptor name.
func (flag *DescriptorFlag) GetValue() string {
	return flag.descriptorName
}

/**************************** Derived Value Flag ******************************/

// TODO: remove the derived flag after migrating to execution engine

// DerivedFlag is used to mark derived values.
type DerivedFlag struct {
	baseKey string
}

// GetIndex returns 5.
func (flag *DerivedFlag) GetIndex() int {
	return DerivedFlagIndex
}

// GetName return name of the Derived flag.
func (flag *DerivedFlag) GetName() string {
	return DerivedFlagName
}

// GetValue returns the key of the base value from which the given derived value
// is derived from (directly or transitively).
func (flag *DerivedFlag) GetValue() string {
	return flag.baseKey
}

/**************************** Value Source Flag *******************************/

// ValueSourceFlag stores all the value source and their pre-merge values.
type ValueSourceFlag struct {
	Sources []exec.ValueWithSource  // ordered by DerivedFrom

	// pre-computed on update
	isDerivedOnly bool // without empty DerivedFrom
	isObtained    bool // without FromNB
}

// GetIndex returns 5.
func (flag *ValueSourceFlag) GetIndex() int {
	return ValueSourceFlagIndex
}

// GetName return name of the Value-Source flag.
func (flag *ValueSourceFlag) GetName() string {
	return ValueStateFlagName
}

// GetValue describes value sources (keys only).
func (flag *ValueSourceFlag) GetValue() string {
	str := "{"
	for i, source := range flag.Sources {
		if source.DerivedFrom == "" {
			str += source.Origin.String()
		} else {
			str += source.DerivedFrom
		}
		if i < len(flag.Sources)-1 {
			str += ", "
		}
	}
	str += "}"
	return str
}

// WithSource add value source into the flag.
// If <inPlace> is disabled, the flag is first copied so that the original remains
// unchanged.
// If <flag> is nil a new flag is created with the single source.
func (flag *ValueSourceFlag) WithSource(source exec.ValueWithSource, inPlace bool) *ValueSourceFlag {

	// TODO
	return nil
}

// WithoutSource removed value source from the flag.
// If <inPlace> is disabled, the flag is first copied so that the original remains
// unchanged.
// If this is the only source, the function returns nil.
func (flag *ValueSourceFlag) WithoutSource(source exec.ValueSource, inPlace bool) *ValueSourceFlag {

	// TODO
	return nil
}

// GetSources returns the value sources, also handling nil flag.
func (flag *ValueSourceFlag) GetSources() []exec.ValueWithSource {
	if flag == nil {
		return nil
	}
	return flag.Sources
}

// GetSourceValue returns the (pre-merge) value associated with the given source.
func (flag *ValueSourceFlag) GetSourceValue(source exec.ValueSource) proto.Message {
	if flag == nil {
		return nil
	}
	// TODO
	return nil
}

// IsDerivedOnly returns true if all the value sources are derived values.
func (flag *ValueSourceFlag) IsDerivedOnly() bool {
	if flag == nil {
		return false
	}
	return flag.isDerivedOnly
}

// IsObtained returns true if all the value sources are from SB.
func (flag *ValueSourceFlag) IsObtained() bool {
	if flag == nil {
		return false
	}
	return flag.isObtained
}