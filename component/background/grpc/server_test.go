package grpc

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/Clink-n-Clank/Brokkr/component/background"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(
		NewServerOptionsBuilder().AddAdditionalGrpcOptions(grpc.InitialConnWindowSize(0)),
	)

	go func() {
		if err := srv.OnStart(ctx); err != nil {
			assert.NoError(t, err, "Unexpected error OnStart gRPC server")
			panic(err)
		}
	}()

	time.Sleep(time.Second)

	err := srv.OnStop(ctx)
	assert.NoError(t, err, "Unexpected error OnStop gRPC server")
	assert.Equal(t, processName, srv.GetName())
	assert.Equal(t, background.TaskSeverityMajor, srv.GetSeverity())
}

func TestListener(t *testing.T) {
	lis := &net.TCPListener{}
	s := NewServer(NewServerOptionsBuilder().AddListener(lis))

	if !reflect.DeepEqual(lis, s.listener) {
		t.Errorf("expect %v, got %v", lis, s.listener)
	}
}

func TestNetwork(t *testing.T) {
	v := "abc"
	s := NewServer(NewServerOptionsBuilder().AddNetwork(v))

	if !reflect.DeepEqual(v, s.network) {
		t.Errorf("expect %s, got %s", v, s.network)
	}
}

func TestAddress(t *testing.T) {
	v := "abc"
	s := NewServer(NewServerOptionsBuilder().AddAddress(v))

	if !reflect.DeepEqual(v, s.address) {
		t.Errorf("expect %s, got %s", v, s.address)
	}
}

func TestSetServicesChecks(t *testing.T) {
	const testingDependedService = "unit"
	healthChecks := map[string]func() grpc_health_v1.HealthCheckResponse_ServingStatus{
		testingDependedService: func() grpc_health_v1.HealthCheckResponse_ServingStatus {
			return grpc_health_v1.HealthCheckResponse_SERVING
		},
	}

	s := NewServer(NewServerOptionsBuilder().AddServicesHealthChecks(healthChecks))
	s.health.Resume()
	r, rErr := s.health.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: testingDependedService})
	assert.NoError(t, rErr)
	assert.Equal(t, r.GetStatus(), grpc_health_v1.HealthCheckResponse_SERVING)

	s.health.Shutdown()
	r, rErr = s.health.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: testingDependedService})
	assert.NoError(t, rErr)
	assert.Equal(t, r.GetStatus(), grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

func TestSetMiddlewares(t *testing.T) {
	b := NewServerOptionsBuilder().
		AddShutdownTimeout(time.Second).
		AddGrpcUnaryInterceptors(grpc_auth.UnaryServerInterceptor(someOtherGRPCUnaryInterceptorForAuth())).
		AddCustomUnaryMiddlewares("*", customMiddleware)

	// s.Server.srvOpts will have 2 new chainUnaryInts, I need Test client to execute integration test to check gRPC middleware
	// Right now you need to call debug to look up this value.
	s := NewServer(b)

	assert.Equal(t, time.Second, s.timeout)
}

func someOtherGRPCUnaryInterceptorForAuth() grpc_auth.AuthFunc {
	return func(ctx context.Context) (context.Context, error) {
		return ctx, nil
	}
}

func customMiddleware(handler RequestHandler) RequestHandler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return handler(ctx, req)
	}
}
