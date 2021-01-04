package models

// modeledKey is a memory-wise minimalistic representation of a key with a model.
type modeledKey struct {
	model      ModelMeta
	itemName   string
	attrSuffix string
}

// key implements KeyV2 interface for keys with models.
type key struct {
	modeledKey
	fullKey string // cache; empty until String() is called (try to avoid it)
}


// stringKey is KeyV2 implementation used with legacy KV descriptors (version 1).
type stringKey struct {
	key string
}

func StringKey(key string) KeyV2 {
	return stringKey{key: key}
}
