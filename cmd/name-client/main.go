package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"time"

	servicepb "github.com/fillmore-labs/name-service/api/fillmore-labs/name-service/v1alpha1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const timeout = 10 * time.Second

type NameServiceFunction struct {
	conn *grpc.ClientConn
}

func randomName(base string) (string, error) {
	id, err := rand.Int(rand.Reader, big.NewInt(1_000)) //nolint:gomnd
	if err != nil {
		return "", fmt.Errorf("can't generate random ID: %w", err)
	}

	return fmt.Sprintf("%s %v", base, id), nil
}

func (f *NameServiceFunction) Sample(ctx context.Context) error {
	givenName, err := randomName("First")
	if err != nil {
		return fmt.Errorf("can't generate random name: %w", err)
	}
	surname, err := randomName("Last")
	if err != nil {
		return fmt.Errorf("can't generate random name: %w", err)
	}

	client := servicepb.NewNameServiceClient(f.conn)

	_, err = client.AddName(ctx, &servicepb.AddNameRequest{GivenName: givenName, Surname: &surname})
	if err != nil {
		return fmt.Errorf("can't add name: %w", err)
	}

	_, err = client.AddName(ctx, &servicepb.AddNameRequest{GivenName: givenName, Surname: nil})
	if err != nil {
		return fmt.Errorf("can't add name: %w", err)
	}

	stream, err := client.ListNames(ctx, &servicepb.ListNamesRequest{})
	if err != nil {
		return fmt.Errorf("can't list names: %w", err)
	}

	for {
		name, err := stream.Recv()
		if err == io.EOF { //nolint:errorlint
			break
		}

		if err != nil {
			return fmt.Errorf("can't list names: %w", err)
		}

		if surname := name.Surname; surname == nil {
			fmt.Println(name.GetGivenName())
		} else {
			fmt.Printf("%s %s\n", name.GetGivenName(), *surname)
		}
	}

	fmt.Println("Done.")

	return nil
}

func New() *NameServiceFunction {
	return &NameServiceFunction{}
}

func (f *NameServiceFunction) Start(_ context.Context) error {
	serverAddr := os.Getenv("NAME_SERVICE")

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial(serverAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial: %w", err)
	}

	f.conn = conn

	return nil
}

func (f *NameServiceFunction) Stop(_ context.Context) error {
	if f.conn == nil {
		return nil
	}

	return f.conn.Close()
}

func main() {
	f := New()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := f.Start(ctx); err != nil {
		log.Fatal(err) //nolint:gocritic
	}

	err := f.Sample(ctx)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			log.Fatal(s.Message())
		}
		log.Fatal(err)
	}

	if err := f.Stop(ctx); err != nil {
		log.Fatal(err)
	}
}
