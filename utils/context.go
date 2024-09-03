package utils

import (
	"context"
	"fmt"
)

type ContextKey string

const AppMode ContextKey = "app_mode"

func AddToContext(ctx context.Context, value interface{}, key ContextKey) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetFromContext[T any](ctx context.Context, key ContextKey) (result T, err error) {
	value := ctx.Value(key)
	if value == nil {
		return result, fmt.Errorf("%s does not exist in context", key)
	}
	result, ok := value.(T)
	if !ok {
		return result, fmt.Errorf("failed to cast %s to %T", key, result)
	}
	return result, nil
}

func GetAppMode(ctx context.Context) string {
	app_mode, err := GetFromContext[string](ctx, AppMode)
	if err != nil {
		return "development"
	}
	return app_mode
}
