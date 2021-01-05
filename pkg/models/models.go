//  Copyright (c) 2018 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package models

import (
	"path"
	"reflect"
	"strings"

	"github.com/go-errors/errors"
	"github.com/golang/protobuf/proto"
	protoV2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/dynamicpb"
)

// Register registers model in DefaultRegistry.
func Register(pb proto.Message, spec Spec, opts ...ModelOption) ProtoModel {
	model, err := DefaultRegistry.Register(pb, spec, opts...)
	if err != nil {
		panic(err)
	}
	return model
}

// RegisterRemote registers remotely known model in given RemoteRegistry
func RegisterRemote(remoteModel *ModelInfo, remoteRegistry *RemoteRegistry) {
	remoteRegistry.Register(remoteModel, ToSpec(remoteModel.Spec))
}

// RegisteredModels returns models registered in the DefaultRegistry.
func RegisteredModels() []ProtoModel {
	return DefaultRegistry.RegisteredModels()
}

// GetModel returns registered model for given model name.
func GetModel(name string) (ProtoModel, error) {
	return GetModelFromRegistry(name, DefaultRegistry)
}

// GetModel returns registered model in given registry for given model name.
func GetModelFromRegistry(name string, modelRegistry Registry) (ProtoModel, error) {
	return modelRegistry.GetModel(name)
}

// GetModelFor returns model registered in DefaultRegistry for given proto message.
func GetModelFor(x proto.Message) (ProtoModel, error) {
	return GetModelFromRegistryFor(x, DefaultRegistry)
}

// GetModelFromRegistryFor returns model registered in modelRegistry for given proto message
func GetModelFromRegistryFor(x proto.Message, modelRegistry Registry) (ProtoModel, error) {
	return modelRegistry.GetModelFor(x)
}

// GetModelForKey returns model registered in DefaultRegistry which matches key.
func GetModelForKey(key KeyV2) (ModelMeta, error) {
	return DefaultRegistry.GetModelForKey(key)
}

// Key is a helper for the GetKey which panics on errors.
func Key(x proto.Message) KeyV2 {
	key, err := GetKey(x)
	if err != nil {
		panic(err)
	}
	return key
}

// Name is a helper for the GetName which panics on errors.
func Name(x proto.Message) string {
	name, err := GetName(x)
	if err != nil {
		panic(err)
	}
	return name
}

// GetKey returns complete key for given model,
// including key prefix defined by model specification.
// It returns error if given model is not registered.
func GetKey(x proto.Message) (KeyV2, error) {
	return GetKeyUsingModelRegistry(x, DefaultRegistry)
}

// GetKey returns complete key for given model from given model registry,
// including key prefix defined by model specification.
// It returns error if given model is not registered.
func GetKeyUsingModelRegistry(message proto.Message, modelRegistry Registry) (KeyV2, error) {
	// find model for message
	model, err := GetModelFromRegistryFor(message, modelRegistry)
	if err != nil {
		return "", errors.Errorf("can't find known model "+
			"for message (while getting key for model) due to: %v (message = %+v)", err, message)
	}

	// compute Item.ID.Name
	name, err := model.InstanceName(message)
	if err != nil {
		return "", errors.Errorf("can't compute model instance name due to: %v (message %+v)", err, message)
	}

	key := path.Join(model.KeyPrefix(), name)
	return key, nil
}

// GetName returns instance name for given model.
// It returns error if given model is not registered.
func GetName(x proto.Message) (string, error) {
	model, err := GetModelFor(x)
	if err != nil {
		return "", err
	}
	name, err := model.InstanceName(x)
	if err != nil {
		return "", err
	}
	return name, nil
}

// AttributeKey is a helper for the GetAttributeKey which panics on errors.
func AttributeKey(attrModelName string, itemVal, attrVal proto.Message) KeyV2 {
	// TODO
	return nil
}

// GetAttributeKey returns complete key for the given attribute instance,.
// It returns error if the given item/attribute model is not registered.
func GetAttributeKey(attrModelName string, itemVal, attrVal proto.Message) (KeyV2, error) {
	// TODO
	return nil, nil
}

// keyPrefix computes correct key prefix from given model. It
// handles correctly the case when name suffix of the key is empty
// (no template name -> key prefix does not end with "/")
func keyPrefix(modelSpec Spec, hasTemplateName bool) string {
	keyPrefix := modelSpec.KeyPrefix()
	if !hasTemplateName {
		keyPrefix = strings.TrimSuffix(keyPrefix, "/")
	}
	return keyPrefix
}

// dynamicMessageToGeneratedMessage converts proto dynamic message to corresponding generated proto message
// (identified by go type).
// This conversion method should help handling dynamic proto messages in mostly protoc-generated proto message
// oriented codebase (i.e. help for type conversions to named, help handle missing data fields as seen
// in generated proto messages,...)
func dynamicMessageToGeneratedMessage(dynamicMessage *dynamicpb.Message,
	goTypeOfGeneratedMessage reflect.Type) (proto.Message, error) {

	// create empty proto message of the same type as it was used for registration
	var registeredGoType interface{}
	if goTypeOfGeneratedMessage.Kind() == reflect.Ptr {
		registeredGoType = reflect.New(goTypeOfGeneratedMessage.Elem()).Interface()
	} else {
		registeredGoType = reflect.Zero(goTypeOfGeneratedMessage).Interface()
	}
	message, isProtoV1 := registeredGoType.(proto.Message)
	if !isProtoV1 {
		messageV2, isProtoV2 := registeredGoType.(protoV2.Message)
		if !isProtoV2 {
			return nil, errors.Errorf("registered go type(%T) is not proto.Message", registeredGoType)
		}
		message = proto.MessageV1(messageV2)
	}

	// fill empty proto message with data from its dynamic proto message counterpart
	// (alternative approach to this is marshalling dynamicMessage to json and unmarshalling it back to message)
	proto.Merge(message, dynamicMessage)

	return message, nil
}
