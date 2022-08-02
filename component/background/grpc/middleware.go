package grpc

import (
	"context"
	"sort"
	"strings"
)

// Abstract middleware functionality based on chain of responsibility

const (
	// globalMiddlewareFilterEmpty will be used if we have middlewares like * (always for all routes)
	globalMiddlewareFilterEmpty = ""
)

// RequestHandler will be invoked in the gRPC Middleware
type RequestHandler func(ctx context.Context, req interface{}) (interface{}, error)

// Middleware to do some actions between gRPC requests
type Middleware func(RequestHandler) RequestHandler

// MiddlewareComposer keeps middlewares and keeps it sorted by filter
type MiddlewareComposer struct {
	mws               []Middleware
	routesUnderFilter []string
	routes            map[string][]Middleware
}

// NewMiddlewareComposer instance
func NewMiddlewareComposer() *MiddlewareComposer {
	return &MiddlewareComposer{
		mws:               make([]Middleware, 0),
		routesUnderFilter: make([]string, 0),
		routes:            map[string][]Middleware{},
	}
}

// Register middleware with filter
func (mc *MiddlewareComposer) Register(filter string, mw ...Middleware) {
	if filter == globalMiddlewareFilterEmpty {
		filter = "*"
	}

	if strings.HasSuffix(filter, "*") {
		filter = strings.TrimSuffix(filter, "*")
		mc.routesUnderFilter = append(mc.routesUnderFilter, filter)

		sort.Slice(mc.routesUnderFilter, func(f1, f2 int) bool {
			return mc.routesUnderFilter[f1] > mc.routesUnderFilter[f2]
		})
	}

	mc.routes[filter] = append(mc.routes[filter], mw...)
}

// Search middleware that related to filtered request
func (mc *MiddlewareComposer) Search(requestPath string) []Middleware {
	ms := make([]Middleware, 0)
	if globalMiddlewares, exist := mc.routes[globalMiddlewareFilterEmpty]; exist {
		ms = append(ms, globalMiddlewares...)
	}

	if requestPath == globalMiddlewareFilterEmpty {
		return ms
	}

	if neededMiddleware, ok := mc.routes[requestPath]; ok {
		return append(ms, neededMiddleware...)
	}

	for _, filtered := range mc.routesUnderFilter {
		if filtered == globalMiddlewareFilterEmpty { // Included by default
			continue
		}

		if strings.HasPrefix(requestPath, filtered) {
			return append(ms, mc.routes[filtered]...)
		}
	}

	return ms
}

// PassToNext delegate to next middleware to execute
func (mc *MiddlewareComposer) PassToNext(m ...Middleware) Middleware {
	return func(requestHandler RequestHandler) RequestHandler {
		for i := len(m) - 1; i >= 0; i-- {
			requestHandler = m[i](requestHandler)
		}

		return requestHandler
	}
}
