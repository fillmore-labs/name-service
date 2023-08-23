package server

import (
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const system = "" // empty string represents the health of the system

func NewHealth() healthpb.HealthServer {
	srv := health.NewServer()
	srv.SetServingStatus(system, healthpb.HealthCheckResponse_SERVING)

	return srv
}
