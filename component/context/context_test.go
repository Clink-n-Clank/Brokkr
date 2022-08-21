package context

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExtendedContextWithMetadata(t *testing.T) {
	var extContext context.Context
	var metaData any
	var isMetaExist bool

	extContext = ExtendedContextWithMetadata[string](context.Background(), "foo", "bar")
	metaData, isMetaExist = GetContextMetadata[string](extContext, "foo")
	assert.True(t, isMetaExist)
	assert.True(t, "string" == fmt.Sprintf("%T", metaData))
	assert.True(t, "bar" == metaData)

	type customStructKey struct{}
	type customStructData struct {
		MyVal string
	}
	extContext = ExtendedContextWithMetadata[customStructData](context.Background(), customStructKey{}, customStructData{"unit"})
	metaData, isMetaExist = GetContextMetadata[customStructData](extContext, customStructKey{})
	assert.True(t, isMetaExist)
	assert.IsTypef(t, customStructData{}, metaData, "expected to be same type of struct")
	assert.True(t, "unit" == metaData.(customStructData).MyVal)
}

func TestForkContextAsNewWithTimeout(t *testing.T) {
	baseCtx := context.Background()
	ctxFork, ctxForkCancel := ForkContextAsNewWithTimeout(baseCtx, time.Minute)
	assert.NotSame(t, baseCtx, ctxFork)
	assert.NotNil(t, ctxForkCancel)
}
