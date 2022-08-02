package grpc

import (
	"context"
	"sort"
	"strings"
)

// Abstract middleware functionality based on chain of responsibility

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
	if filter == "" {
		filter = "*"
	}

	if strings.HasSuffix(filter, "*") {
		filter = strings.TrimSuffix(filter, "*")
		mc.routesUnderFilter = append(mc.routesUnderFilter, filter)

		sort.Slice(mc.routesUnderFilter, func(f1, f2 int) bool {
			return mc.routesUnderFilter[f1] < mc.routesUnderFilter[f2]
		})
	}

	mc.routes[filter] = append(mc.routes[filter], mw...)
}

// Search middleware that related to filtered request
func (mc *MiddlewareComposer) Search(requestPath string) []Middleware {
	lookupChain := map[string][]Middleware{}

	for _, filtered := range mc.routesUnderFilter {
		if strings.HasPrefix(requestPath, filtered) {
			lookupChain[filtered] = mc.routes[filtered]
		}
	}

	ms := make([]Middleware, 0)
	for _, m := range lookupChain {
		ms = append(ms, m...)
	}

	return ms
}

// PassToNext delegate to next middleware to execute
func (mc *MiddlewareComposer) PassToNext(m ...Middleware) Middleware {
	return func(next RequestHandler) RequestHandler {
		for i := len(m) - 1; i >= 0; i-- {
			next = m[i](next)
		}

		return next
	}
}
