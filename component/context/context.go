package context

import (
	"context"
	"time"
)

// ExtendedContextWithMetadata add metadata to base context
func ExtendedContextWithMetadata[metaType any](baseCtx context.Context, metaKey any, metadata metaType) context.Context {
	return context.WithValue(baseCtx, metaKey, metadata)
}

// GetContextMetadata will try get form context.Context metadata
func GetContextMetadata[metaType any](baseCtx context.Context, metaKey any) (metaData metaType, isExist bool) {
	metaData, isExist = baseCtx.Value(metaKey).(metaType)

	return
}

// ForkContextAsNewWithTimeout will take base context and create new from it
func ForkContextAsNewWithTimeout(baseCtx context.Context, newTimeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(baseCtx, newTimeout)
}
