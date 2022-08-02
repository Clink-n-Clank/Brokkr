package grpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/Clink-n-Clank/Brokkr/component/background"
)

// BackgroundServer wrapper
type BackgroundServer struct {
	*grpc.Server
	opts               []grpc.ServerOption
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	health             *health.Server
	middlewareComposer *MiddlewareComposer

	network string
	address string
	timeout time.Duration

	listener    net.Listener
	listenerErr error

	// dependedServicesCheck has as string - service name and function that returns state
	dependedServicesCheck map[string]func() grpc_health_v1.HealthCheckResponse_ServingStatus
}

const (
	processName = "gRPC Server"
	netProtocol = "tcp"
	netAddress  = ":0"
)

// NewServer instance
func NewServer(builder *ServerOptionsBuilder) *BackgroundServer {
	// Set defaults for new process wrapper
	serv := &BackgroundServer{
		network:            netProtocol,
		address:            netAddress,
		timeout:            30 * time.Second,
		health:             health.NewServer(),
		middlewareComposer: NewMiddlewareComposer(),
	}

	// Load additional grpc server options
	for _, o := range builder.Build() {
		o(serv)
	}

	// TODO add stream middleware support
	// Load middlewares and interceptors
	serverUnaryInterceptor := []grpc.UnaryServerInterceptor{serv.unaryServerInterceptorForMiddleware()}
	if len(serv.unaryInterceptors) > 0 {
		serverUnaryInterceptor = append(serverUnaryInterceptor, serv.unaryInterceptors...)
	}

	serv.opts = append(
		serv.opts,
		[]grpc.ServerOption{
			grpc.ChainUnaryInterceptor(serverUnaryInterceptor...),
			grpc.ChainStreamInterceptor(serv.streamInterceptors...),
		}...,
	)

	// Create and run gRPC server
	serv.Server = grpc.NewServer(serv.opts...)
	serv.listenerErr = serv.listen()

	// Add internal sub-services to gRPC server register
	for serviceName, check := range serv.dependedServicesCheck {
		serv.health.SetServingStatus(serviceName, check())
	}

	grpc_health_v1.RegisterHealthServer(serv.Server, serv.health)

	return serv
}

// GetName of the task
func (s *BackgroundServer) GetName() string {
	return processName
}

// GetSeverity of the task
func (s *BackgroundServer) GetSeverity() background.ProcessSeverity {
	return background.TaskSeverityMajor
}

// OnStart event to be called when main loop will be started
func (s *BackgroundServer) OnStart(_ context.Context) error {
	if s.listenerErr != nil {
		return s.listenerErr
	}

	s.health.Resume()

	return s.Serve(s.listener)
}

// OnStop event to be called when main loop will be started
func (s *BackgroundServer) OnStop(_ context.Context) error {
	s.health.Shutdown()
	s.GracefulStop()

	return nil
}

// Listen network traffic for service handling
func (s *BackgroundServer) listen() error {
	if s.listener == nil {

		lis, err := net.Listen(s.network, s.address)
		if err != nil {
			return err
		}

		s.listener = lis
	}

	return nil
}

// unaryServerInterceptorForMiddleware managing gRPC request interception to delegate it to Middleware
func (s *BackgroundServer) unaryServerInterceptorForMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//
		// Look up for registered middlewares
		//
		defaultRequestHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
			return handler(ctx, req)
		}

		affectedMiddlewares := s.middlewareComposer.Search(info.FullMethod)
		if len(affectedMiddlewares) > 0 {
			defaultRequestHandler = s.middlewareComposer.PassToNext(affectedMiddlewares...)(defaultRequestHandler)
		}

		return defaultRequestHandler(ctx, req)
	}
}

// TODO add stream middleware support
