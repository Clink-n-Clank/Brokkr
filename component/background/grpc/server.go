package grpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/Clink-n-Clank/Brokkr/component/background"
)

type (
	// BackgroundServer wrapper
	BackgroundServer struct {
		*grpc.Server

		opts []grpc.ServerOption

		network string
		address string
		timeout time.Duration

		listener    net.Listener
		listenerErr error
	}

	// Options sets options such as credentials, keepalive parameters, etc.
	Options func(o *BackgroundServer)
)

const (
	processName = "gRPC Server"
	netProtocol = "tcp"
	netAddress  = ":0"
)

// SetNetwork custom value
func SetNetwork(n string) Options {
	return func(s *BackgroundServer) {
		s.network = n
	}
}

// SetListener custom value
func SetListener(l net.Listener) Options {
	return func(s *BackgroundServer) {
		s.listener = l
	}
}

// SetAddress  custom value
func SetAddress(a string) Options {
	return func(s *BackgroundServer) {
		s.address = a
	}
}

// SetTimeout custom value
func SetTimeout(t time.Duration) Options {
	return func(s *BackgroundServer) {
		s.timeout = t
	}
}

// AddOptions for gRPC server
func AddOptions(opts ...grpc.ServerOption) Options {
	return func(s *BackgroundServer) {
		s.opts = opts
	}
}

// NewServer instance
func NewServer(opts ...Options) *BackgroundServer {
	// Set defaults for new process wrapper
	serv := &BackgroundServer{
		network: netProtocol,
		address: netAddress,
		timeout: 30 * time.Second,
	}

	// Load additional grpc server options
	for _, o := range opts {
		o(serv)
	}

	serv.Server = grpc.NewServer(serv.opts...)
	serv.listenerErr = serv.listen()

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

	// TODO Add logger as dependency -> BackgroundServer start

	return s.Serve(s.listener)
}

// OnStop event to be called when main loop will be started
func (s *BackgroundServer) OnStop(_ context.Context) error {
	s.GracefulStop()

	// TODO Add logger as dependency -> BackgroundServer stop

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
