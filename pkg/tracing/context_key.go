package tracing

import (
	"context"

	"github.com/google/uuid"
)

type ContextKey string

const (
	ContextTraceKey ContextKey = "TraceID"
)

type ContextTrace struct {
	key ContextKey
}

func NewContextTrace() ContextTrace {
	return ContextTrace{key: ContextTraceKey}
}

func (c ContextTrace) Getkey() ContextKey {
	return c.key
}

func (c ContextTrace) GetValueFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value(c.Getkey()).(string); ok && id != "" {
		return id
	}
	return ""
}

func (c ContextTrace) GenerateID() string {
	id := uuid.New().String()
	return id
}
