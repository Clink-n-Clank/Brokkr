package grpc

import (
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Options sets options such as credentials, keepalive parameters, etc.
type Options func(o *BackgroundServer)

// ServerOptionsBuilder sets options such as credentials, keepalive parameters, etc, related to gRPC server
type ServerOptionsBuilder struct {
	srvOpts []Options
}

// NewServerOptionsBuilder for gRPC server configuration
func NewServerOptionsBuilder() *ServerOptionsBuilder {
	return &ServerOptionsBuilder{srvOpts: make([]Options, 0)}
}

// AddNetwork type that will be used in gRPC server
func (b *ServerOptionsBuilder) AddNetwork(n string) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.network = n })
	return b
}

// AddAddress that will be used in gRPC server endpoint
func (b *ServerOptionsBuilder) AddAddress(a string) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.address = a })
	return b
}

// AddListener custom value
func (b *ServerOptionsBuilder) AddListener(l net.Listener) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.listener = l })
	return b
}

// AddShutdownTimeout for gRPC server
func (b *ServerOptionsBuilder) AddShutdownTimeout(t time.Duration) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.timeout = t })
	return b
}

// AddGrpcUnaryInterceptors intercept the execution of a unary RPC on the server
func (b *ServerOptionsBuilder) AddGrpcUnaryInterceptors(unaryInter ...grpc.UnaryServerInterceptor) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.unaryInterceptors = unaryInter })
	return b
}

// AddGrpcStreamInterceptors intercept the execution of a streaming RPC on the server
func (b *ServerOptionsBuilder) AddGrpcStreamInterceptors(streamInter ...grpc.StreamServerInterceptor) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.streamInterceptors = streamInter })
	return b
}

// AddCustomUnaryMiddlewares that have first priority to intercepted request in the middlewares and forwards it to gRPC if needed
// filter used for calling middleware for example:
// - /myapp.v1.MyAppAPI/*                     - Middleware will be executed for all endpoints under "/myapp.v1.MyAppAPI"
// - /myapp.v1.MyAppAPI/OnlyThatEndpoint      - Middleware will be executed only for "OnlyThatEndpoint"
func (b *ServerOptionsBuilder) AddCustomUnaryMiddlewares(filter string, mwList ...Middleware) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.middlewareComposer.Register(filter, mwList...) })
	return b
}

// AddServicesHealthChecks to verify if gRPC working correctly as health checks
func (b *ServerOptionsBuilder) AddServicesHealthChecks(srv map[string]func() grpc_health_v1.HealthCheckResponse_ServingStatus) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.dependedServicesCheck = srv })
	return b
}

// AddAdditionalGrpcOptions that's needed for gRPC server
func (b *ServerOptionsBuilder) AddAdditionalGrpcOptions(grpcOpts ...grpc.ServerOption) *ServerOptionsBuilder {
	b.srvOpts = append(b.srvOpts, func(s *BackgroundServer) { s.opts = grpcOpts })
	return b
}

// Build will make sure that all needed options prepared for server
func (b *ServerOptionsBuilder) Build() []Options {
	return b.srvOpts
}
