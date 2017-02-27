package wrpctx

import "context"

type keyType string
type ctxMap map[keyType]interface{}

const mapKey = "wrpctx"

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

func New(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyType(mapKey), make(ctxMap))
}

func NewWithValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, keyType(key), value)
}

func GetCtxValue(ctx context.Context, key string) interface{} {
	return ctx.Value(keyType(key))
}

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
