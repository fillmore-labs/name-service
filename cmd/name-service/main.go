package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	servicepb "github.com/fillmore-labs/name-service/api/fillmore-labs/name-service/v1alpha1"
	"github.com/fillmore-labs/name-service/database"
	"github.com/fillmore-labs/name-service/server"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const (
	timeout     = 10 * time.Second
	defaultPort = "8080"
)

func main() {
	fmt.Println("Starting project service")

	if err := runAll(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done")
}

func runAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db, err := database.NewDatabase(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Disconnect(ctx)
	}()

	lis, err := NewListener()
	if err != nil {
		return err
	}
	defer func() {
		_ = lis.Close()
	}()

	grpcServer := NewServer(db)

	return RunServer(grpcServer, lis)
}

func NewListener() (net.Listener, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	fmt.Printf("Listening on port %s\n", port)

	return lis, nil
}

func NewServer(db database.Database) *grpc.Server {
	srv := server.NewServer(db)
	health := server.NewHealth()

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)

	servicepb.RegisterNameServiceServer(grpcServer, srv)
	healthpb.RegisterHealthServer(grpcServer, health)
	reflection.Register(grpcServer)

	return grpcServer
}

func RunServer(grpcServer *grpc.Server, lis net.Listener) error {
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm

		fmt.Println("Terminating")

		grpcServer.GracefulStop()
	}()

	return grpcServer.Serve(lis)
}
