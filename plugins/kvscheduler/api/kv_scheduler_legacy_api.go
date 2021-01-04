package api

import (
	"github.com/golang/protobuf/proto"
)

// KeySelectorV1 is a legacy definition of Key Selector.
type KeySelectorV1 = KeySelector

// KeySelector is used to filter keys.
type KeySelector func(key string) bool

// KeyValuePairV1 is a legacy definition of Key Value Pair.
type KeyValuePairV1 = KeyValuePair

// KeyValuePair groups key with value.
type KeyValuePair struct {
	// Key identifies value.
	Key string

	// Value may represent some object, action or property.
	//
	// Value can be created either via northbound transaction (NB-value,
	// ValueOrigin = FromNB) or pushed (as already created) through SB notification
	// (SB-value, ValueOrigin = FromSB). Values from NB take priority as they
	// overwrite existing SB values (via Modify operation), whereas notifications
	// for existing NB values are ignored. For values retrieved with unknown
	// origin the scheduler reviews the value's history to determine where it came
	// from.
	//
	// For descriptors the values are mutable objects - Create, Update and Delete
	// methods should reflect the value content without changing it.
	// To add and maintain extra (runtime) attributes alongside the value, descriptor
	// can use the value metadata.
	Value proto.Message
}

// KVWithMetadataV1 is a legacy definition of KV With Metadata.
type KVWithMetadataV1 = KVWithMetadata

// KVWithMetadata encapsulates key-value pair with metadata and the origin mark.
type KVWithMetadata struct {
	Key      string
	Value    proto.Message
	Metadata Metadata
	Origin   ValueOrigin
}

