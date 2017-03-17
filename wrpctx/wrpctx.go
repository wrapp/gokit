package wrpctx

import "context"

type keyType string
type ctxMap map[keyType]interface{}

const mapKey = "wrpctx"

// Set adds a new value in the context against the provided key. This uses a map
// internally to convinently store key values. Using a map internally allows to
// iterate keys and values stored in the context which is not provided by default
// context behaviour. The values can be retrieved by Get fuction.
func Set(ctx context.Context, key string, value interface{}) {
	cm := ctx.Value(keyType(mapKey))

	if cm == nil {
		return
	}

	m, ok := cm.(ctxMap)
	if !ok {
		return
	}

	m[keyType(key)] = value
}

// Get gets the stored key which was set using Set function. If there is no such key nil
// is returned.
func Get(ctx context.Context, key string) interface{} {
	cm := ctx.Value(keyType(mapKey))

	if cm == nil {
		return nil
	}

	m, ok := cm.(ctxMap)
	if !ok {
		return nil
	}

	return m[keyType(key)]
}

// New creates a new context. It sets up an internal map in the provided context.
func New(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyType(mapKey), make(ctxMap))
}

// NewWithValue creates and returns a new context with the provided value set. This does not use
// the internal map which was initialzed in New function. This creates a new value in the
// context which you can do manually by calling WithValue in context.Context package. This
// calue can be retrieved by GetCtxValue or Value function on context.Context. Recommended
// way to get this value is through GetCtxValue.
func NewWithValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, keyType(key), value)
}

// GetCtxValue returns the value which was stored through NewWithValue. This does not use
// internal map but rather uses Value in context.Context.
func GetCtxValue(ctx context.Context, key string) interface{} {
	return ctx.Value(keyType(key))
}

// GetMap returns a copy of internal map. This allows the user of wrapp context to iterate
// and get all keys and values of internal map. Modifying the returned map will not affect
// the internal map.
func GetMap(ctx context.Context) map[string]interface{} {
	newMap := make(map[string]interface{})
	cm := ctx.Value(keyType(mapKey))

	if cm == nil {
		return newMap
	}

	m, ok := cm.(ctxMap)
	if !ok {
		return newMap
	}

	for key, value := range m {
		newMap[string(key)] = value
	}
	return newMap
}
