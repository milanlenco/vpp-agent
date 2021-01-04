package api

import (
	"context"
	"github.com/golang/protobuf/proto"
)

// WithDerivedFrom puts base value into the context for the descriptor of a derived
// value to obtain if needed.
func WithDerivedFrom(ctx context.Context, msg proto.Message) context.Context {
	// TODO
	return nil
}

// DerivedFrom can be used by descriptor of derived value to obtain the base value
// (value it was derived from).
func DerivedFrom(ctx context.Context) (msg proto.Message, isDerived bool) {
	// TODO
	return nil, false
}

// WithKey puts key of the value into the context for descriptor methods to use if needed.
func WithKey(ctx context.Context, key string) context.Context {
	// TODO
	return nil
}

// Key can be used by descriptor methods to obtain the key of the value that is being
// processed.
func Key(ctx context.Context) (key string) {
	// TODO
	return ""
}
