package wrpctx

import "context"

type keyType string
type ctxMap map[keyType]interface{}

const mapKey = "wrpctx"

func Set(ctx context.Context, key string, value interface{}) {
	if cm := ctx.Value(keyType(mapKey)); cm != nil {
		if m, ok := cm.(ctxMap); ok {
			m[keyType(key)] = value
		}
	}
}

func Get(ctx context.Context, key string) interface{} {
	if cm := ctx.Value(keyType(mapKey)); cm != nil {
		if m, ok := cm.(ctxMap); ok {
			return m[keyType(key)]
		}
	}
	return nil
}

func New(ctx context.Context) context.Context {
	return context.WithValue(ctx, keyType(mapKey), make(ctxMap))
}

func NewWithValue(ctx context.Context, key string, value interface{}) context.Context {
	return context.WithValue(ctx, keyType(key), value)
}

func GetCtxValue(ctx context.Context, key string) interface{} {
	return ctx.Value(keyType(mapKey))
}
