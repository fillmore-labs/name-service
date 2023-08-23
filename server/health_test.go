package server_test

import (
	"context"
	"testing"

	"github.com/fillmore-labs/name-service/server"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const system = "" // empty string represents the health of the system

func TestHealth(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.Background()
	health := server.NewHealth()
	// When
	req := &healthpb.HealthCheckRequest{Service: system}
	res, err := health.Check(ctx, req)
	// Then
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if res.GetStatus() != healthpb.HealthCheckResponse_SERVING {
		t.Errorf("Expected SERVING, got: %v", res.GetStatus().String())
	}
}
