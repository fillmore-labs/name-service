package main_test

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	servicepb "github.com/fillmore-labs/name-service/api/fillmore-labs/name-service/v1alpha1"
	main "github.com/fillmore-labs/name-service/cmd/name-service"
	"github.com/fillmore-labs/name-service/database"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const timeout = 10 * time.Second

var errNotTerminated = errors.New("name service: couldn't shut down gracefully")

type testConnection struct {
	db         database.Database
	lis        net.Listener
	grpcServer *grpc.Server
	srvErr     <-chan error
	conn       *grpc.ClientConn
}

func (c *testConnection) Close(ctx context.Context) error {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return err
		}
	}

	c.grpcServer.GracefulStop()

	err, ok := <-c.srvErr
	if !ok {
		return errNotTerminated
	}
	if err != nil {
		return err
	}

	if err := c.lis.Close(); err != nil {
		return err
	}
	if err := c.db.Disconnect(ctx); err != nil { 
		return err
	}

	return nil
}

func (c *testConnection) Stop(_ context.Context) error {
	c.grpcServer.GracefulStop()

	err, ok := <-c.srvErr
	if !ok {
		return errNotTerminated
	}

	return err
}

func newTestConnection(ctx context.Context) (*testConnection, error) {
	db, err := database.NewDatabase(ctx)
	if err != nil {
		return nil, err
	}

	grpcServer := main.NewServer(db)

	lis := bufconn.Listen(1_024 * 1_024)
	srvErr := make(chan error, 1)
	go func() {
		srvErr <- grpcServer.Serve(lis)
	}()

	contextDialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithContextDialer(contextDialer),
	}

	c := &testConnection{db: db, lis: lis, grpcServer: grpcServer, srvErr: srvErr}

	conn, err := grpc.DialContext(ctx, "", opts...)
	if err != nil {
		_ = c.Close(ctx)

		return nil, err
	}
	c.conn = conn

	return c, nil
}

func TestServer(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	c, err := newTestConnection(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err = c.Close(ctx)
		if err != nil {
			t.Logf("couldn't terminate: %v", err)
		}
	}()

	client := servicepb.NewNameServiceClient(c.conn)

	last := "Last"
	_, err = client.AddName(ctx, &servicepb.AddNameRequest{GivenName: "First", Surname: &last})
	if err != nil {
		t.Fatalf("couldn't add name: %v", err)
	}
}
