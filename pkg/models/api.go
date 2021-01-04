// Copyright (c) 2020 Pantheon.tech
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

package models

import (
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
	"go.ligato.io/vpp-agent/v3/proto/ligato/generic"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// KeyV2 is an immutable identifier of a configuration or a notification instance.
// KeyV2 is composed of a prefix, matching all instances of the same type, and a suffix,
// naming the instance. The prefix itself implements KeyV2, returning empty suffix.
//
// Originally, key was declared as string (version 1). However, key operations are
// very frequent and to parse key using functions from the strings packages was
// quite costly, especially in terms of memory allocations and GC overhead.
// Using interface (key version 2) it is possible to provide more efficient
// implementations for key operations, using for example caching or by keeping
// a direct reference to the associated model.
type KeyV2 interface {
	// Returns string representation of the key.
	// The implementation of KeyV2 may cache the returned string (key is immutable)
	// instead of constructing it with each call.
	// It is, however, recommended to use WriteString instead and potentially avoid
	// memory allocation.
	String() string

	// WriteString writes the key as string into w (preferably without temporary memory
	// allocations).
	WriteString(w io.Writer) error

	// Prefix return longest key prefix used with this key (recursively it is possible
	// to discover all prefixes in the descending order of prefix length).
	// Returns empty Key if there are no further prefixes.
	Prefix() KeyV2

	// Returns empty string if this is key-prefix or a reference to a singleton instance.
	Suffix() string

	// Returns true if the key is empty.
	// Method is equivalent to (but depending on the implementation executes much more efficiently):
	//   func (key1 *KeyV2Impl) Empty() bool {
	//       return key1.String() == ""
	//   }
	Empty() bool

	// Returns true if this key has the given prefix.
	// Method is equivalent to (but depending on the implementation executes much more efficiently):
	//   func (key1 *KeyV2Impl) HasPrefix(key2 KeyV2) bool {
	//       prefix := key1
	//       for !prefix.Empty() {
	//           if key2.Equals(prefix) {
	//               return true
	//           }
	//           prefix = prefix.Prefix()
	//       }
	//       return false
	//   }
	HasPrefix(KeyV2) bool

	// Compare returns 0 if thisKey==key2, -1 if thisKey < key2, and +1 if thisKey > key2.
	// Method is equivalent to (but depending on the implementation executes much more efficiently):
	//   func (key1 *KeyV2Impl) Compare(key2 KeyV2) int {
	//       switch {
	//       case key1.String() == key2.String(): return 0
	//       case key1.String() < key2.String(): return -1
	//       }
	//       return 1
	//   }
	// For equality comparison, without the need for ordering, it is preferred to call Equals
	// instead which could potentially save few instructions.
	Compare(KeyV2) int

	// MapKey returns minimalistic representation of KeyV2 applicable for use as a map key,
	// therefore the underlying type must be comparable (i.e. avoid slices, maps and function values).
	// It should be satisfied that:
	//  key1.Compare(key2) == 0  <=>  key1.MapKey() == key2.MapKey()
	//
	// Important: never use different implementations of KeyV2 with the same key because
	// even when Compare() returns 0 (equal), the "==" operator might return false
	// if the underlying types differ.
	MapKey() interface{}
}

// ModelInfo represents model information retrieved using meta service
type ModelInfo struct {
	generic.ModelDetail

	// MessageDescriptor is the proto message descriptor of the message represented by this ModelInfo struct
	MessageDescriptor protoreflect.MessageDescriptor
}

// Registry defines model registry for managing registered models
type Registry interface {
	// GetModel returns registered model for the given model name
	// or error if model is not found.
	GetModel(name string) (ProtoModel, error)

	// GetModelFor returns registered model for the given proto message.
	GetModelFor(x interface{}) (ProtoModel, error)

	// GetModelForKey returns registered model for the given key or error.
	GetModelForKey(key string) (ProtoModel, error)

	// MessageTypeRegistry creates new message type registry from registered proto messages
	MessageTypeRegistry() *protoregistry.Types

	// RegisteredModels returns all registered modules.
	RegisteredModels() []ProtoModel

	// Register registers either a protobuf message known at compile-time together
	// with the given model specification (for LocalRegistry),
	// or a remote model represented by an instance of ModelInfo obtained via KnownModels RPC from MetaService
	// (for RemoteRegistry or also for LocalRegistry but most likely just proxied to a remote agent).
	// If spec.Class is unset, then it defaults to 'config'.
	Register(x interface{}, spec Spec, opts ...ModelOption) (ProtoModel, error)
}

type ModelMeta interface {
	// Stringer is only used for logging purposes.
	fmt.Stringer

	// Name returns the name for the model.
	// It is unique across all models, including attributes and attribute groups.
	Name() string

	// Class is a string value that classifies the modeled items.
	// For example, "config" is used for configuration items.
	// Currently used values are: "config", "notif", "metrics", "attr", "attr-group".
	Class() string

	// Parent returns pointer to a parent model at the given index.
	// Models are organized in a tree hierarchy with parent-child relationship.
	// For example, Attribute is a child of ProtoModel or AttributeGroup.
	// Each ProtoModel is at the root of its own tree with index 0.
	// Attribute or AttributeGroup registered directly on a ProtoModel has index 1,
	// children of top-most AttributeGroup have index 2, and so on.
	// The last valid index is `NestingDepth()-1`, which points to this model
	// (i.e. will return itself).
	Parent(idx int) ModelMeta

	// NestingDepth returns the number of parent models plus one.
	// For example, the method returns 1 for ProtoModels because they are never nested.
	// Attribute will have depth at least 2 or more if it is inside Attribute Groups.
	NestingDepth() int

	// KeyPrefix returns the longest common key prefix of this model items.
	KeyPrefix() KeyV2

	// Registry returns pointer to registry where this model is registered.
	Registry() Registry
}

// ProtoModel represents a registered model that is associated with a Proto Message.
// There can be at most one ProtoModel defined for a given Proto Message.
type ProtoModel interface {
	ModelMeta

	// Spec returns model specification for the model.
	Spec() *Spec

	// ModelDetail returns descriptor for the model.
	ModelDetail() *generic.ModelDetail

	// NewInstance creates new instance value for model type.
	NewInstance() proto.Message

	// ProtoName returns proto message name registered with the model.
	ProtoName() string

	// ProtoFile returns proto file name for the model.
	ProtoFile() string

	// NameTemplate returns name template for the model.
	NameTemplate() string

	// GoType returns go type for the model.
	GoType() string

	// PkgPath returns package import path for the model definition.
	PkgPath() string

	// ParseKey parses the given key and returns item name
	// or returns empty name and valid as false if the key is not valid.
	ParseKey(key KeyV2) (name string, valid bool)

	// IsKeyValid returns true if given key is valid for this model.
	IsKeyValid(key KeyV2) bool

	// StripKeyPrefix returns key with prefix stripped.
	StripKeyPrefix(key KeyV2) string

	// InstanceName computes message name for given proto message using name template (if present).
	InstanceName(x interface{}) (string, error)

	// RegisterAttribute registers new attribute for this item.
	RegisterAttribute(spec *AttrSpec, opts ...AttrOption) Attribute

	// RegisterAttributeGroup register new group of attributes for this item.
	RegisterAttributeGroup(name string) AttributeGroup
}

// Attribute is a field/action assigned to a ProtoModel and potentially nested under Attribute Group(s).
// Just like for ProtoModel, the value of Attribute also has to be an instance of Proto Message.
// The difference, however, is that there does not have to be a distinct Proto Message defined
// for each Attributes. Instead Attributes may use the same Proto Messages, e.g. from the package
// "protobuf/ptypes" with already prepared protos for well-known types.
type Attribute interface {
	ModelMeta

	// Spec returns model specification for the Attribute.
	Spec() *AttrSpec

	// KeySuffixTemplate returns key suffix template for the attribute.
	KeySuffixTemplate() string

	// InstanceKeySuffix computes key suffix for the given attribute value using the template
	// (if present).
	InstanceKeySuffix(attrVal proto.Message) (string, error)

	// ParseKey parses the given key and returns the attribute key suffix as well as the name
	// of the underlying item.
	// Returns empty strings and valid as false if the key is not valid.
	ParseKey(key KeyV2) (itemName, attrKeySuffix string, valid bool)

	// IsKeyValid returns true if given key is valid for this model.
	IsKeyValid(key KeyV2) bool
}

// AttributeGroup groups related attributes under a common key prefix.
type AttributeGroup interface {
	ModelMeta

	// AttrGroupName returns the name of this attribute group.
	AttrGroupName() string

	// RegisterAttributeSubGroup registers new group of attributes
	// as nested under this group.
	RegisterAttributeSubGroup(name string) AttributeGroup

	// RegisterAttribute registers new attribute under this group.
	RegisterAttribute(spec *AttrSpec, opts ...AttrOption) Attribute
}