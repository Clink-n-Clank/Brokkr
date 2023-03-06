package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

var (
	testingPassToMiddleware = 0
)

func TestRegisterAndSearchOfMiddlewares(t *testing.T) {
	mc := NewMiddlewareComposer()
	mc.Register("*", testingPermanentMiddleware)
	mc.Register("route/A/test/*", testingFilteredMiddleware)
	mc.Register("route/A/*", testingFilteredMiddleware)

	affectedMiddlewares := mc.Search("route/B")
	assert.NotEmpty(t, affectedMiddlewares)
	assert.Len(t, affectedMiddlewares, 1)

	assert.Contains(t, mc.routesUnderFilter, "")
	assert.Contains(t, mc.routesUnderFilter, "route/A/test/")
	assert.Contains(t, mc.routesUnderFilter, "route/A/")
}

func TestOrderMiddlewares(t *testing.T) {
	mc := NewMiddlewareComposer()
	mc.Register("*", testingPermanentMiddleware)
	mc.Register("*", testingThisWillLogEachRequestMiddleware(t))
	mc.Register("route/A/*", testingFilteredMiddleware)

	affectedMiddlewares := mc.Search("route/B")
	assert.NotEmpty(t, affectedMiddlewares)
	assert.Len(t, affectedMiddlewares, 2)
}

func TestPassToNext(t *testing.T) {
	testingPassToMiddleware = 0

	reqHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		testingPassToMiddleware++
		return "test", nil
	}

	mc := NewMiddlewareComposer()
	mc.Register("*")
	returned, err := mc.PassToNext(testingPermanentMiddleware, testingFilteredMiddleware)(reqHandler)(context.Background(), "unit")
	assert.Nil(t, err)

	if !reflect.DeepEqual(returned, "test") {
		t.Errorf("request handle must return: `%s` but returned: `%v`", "test", returned)
	}

	if !reflect.DeepEqual(testingPassToMiddleware, 3) {
		t.Errorf("each middleware must iterate value and request handler as well, sum must be: `3` but returned: %v", testingPassToMiddleware)
	}
}

func TestCheckIfAppendedByPathMiddlewares(t *testing.T) {
	// This case must append middleware by:
	// /* - each time
	// /api.route/* - all what is after "route"
	mc := NewMiddlewareComposer()
	mc.Register("*", testingPermanentMiddleware)
	mc.Register("/api.route/*", testingFilteredMiddleware, testingThisWillLogEachRequestMiddleware(t))

	affectedMiddlewares := mc.Search("/api.route/methods/FooBar")
	assert.NotEmpty(t, affectedMiddlewares)
	assert.Len(t, affectedMiddlewares, 3)
}

func TestContextMetadata(t *testing.T) {
	const expectedFullMethod = "FooBar"
	expectedMetaData := metadata.MD{"foo": {"bar"}}
	ctx := context.Background()

	mc := NewMiddlewareComposer()
	extCtx := mc.ExtendContext(ctx, RequestContextMetadata{
		Meta:       expectedMetaData,
		FullMethod: expectedFullMethod,
	})

	mcCtxMeta, isMcCtxMetaExist := GetContextMetadata(extCtx)

	assert.True(t, isMcCtxMetaExist)
	assert.Equal(t, expectedFullMethod, mcCtxMeta.FullMethod)
	assert.Equal(t, expectedMetaData, mcCtxMeta.Meta)
}

func testingFilteredMiddleware(handler RequestHandler) RequestHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		testingPassToMiddleware++
		return handler(ctx, req)
	}
}

func testingPermanentMiddleware(handler RequestHandler) RequestHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		testingPassToMiddleware++
		return handler(ctx, req)
	}
}

func testingThisWillLogEachRequestMiddleware(t *testing.T) Middleware {
	return func(handler RequestHandler) RequestHandler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			t.Logf("testingThisWillLogEachRequestMiddleware -> req:%v\n", req)
			return handler(ctx, req)
		}
	}
}
