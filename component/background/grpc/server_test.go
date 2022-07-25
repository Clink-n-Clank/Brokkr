package grpc

import (
	"context"
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestServer(t *testing.T) {
	ctx := context.Background()
	srv := NewServer(
		AddOptions(grpc.InitialConnWindowSize(0)),
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
}

func TestOptions(t *testing.T) {
	s := &BackgroundServer{}
	v := []grpc.ServerOption{
		grpc.EmptyServerOption{},
	}

	AddOptions(v...)(s)

	if !reflect.DeepEqual(v, s.opts) {
		t.Errorf("expect %v, got %v", v, s.opts)
	}
}

func TestListener(t *testing.T) {
	lis := &net.TCPListener{}

	s := &BackgroundServer{}
	SetListener(lis)(s)

	if !reflect.DeepEqual(lis, s.listener) {
		t.Errorf("expect %v, got %v", lis, s.listener)
	}
}

func TestNetwork(t *testing.T) {
	s := &BackgroundServer{}
	v := "abc"

	SetNetwork(v)(s)

	if !reflect.DeepEqual(v, s.network) {
		t.Errorf("expect %s, got %s", v, s.network)
	}
}

func TestAddress(t *testing.T) {
	s := &BackgroundServer{}
	v := "abc"

	SetAddress(v)(s)

	if !reflect.DeepEqual(v, s.address) {
		t.Errorf("expect %s, got %s", v, s.address)
	}
}

func TestSetServicesChecks(t *testing.T) {
	const testingDependedService = "unit"
	opts := []Options{
		SetServicesChecks(map[string]func() grpc_health_v1.HealthCheckResponse_ServingStatus{
			testingDependedService: func() grpc_health_v1.HealthCheckResponse_ServingStatus {
				return grpc_health_v1.HealthCheckResponse_SERVING
			},
		}),
	}

	s := NewServer(opts...)

	s.health.Resume()
	r, rErr := s.health.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: testingDependedService})
	assert.NoError(t, rErr)
	assert.Equal(t, r.GetStatus(), grpc_health_v1.HealthCheckResponse_SERVING)

	s.health.Shutdown()
	r, rErr = s.health.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: testingDependedService})
	assert.NoError(t, rErr)
	assert.Equal(t, r.GetStatus(), grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}
