package server_test

import (
	"context"
	"errors"
	"testing"

	servicepb "github.com/fillmore-labs/name-service/api/fillmore-labs/name-service/v1alpha1"
	"github.com/fillmore-labs/name-service/database"
	"github.com/fillmore-labs/name-service/server"
)

func TestRegisterProject(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.Background()
	server := server.NewServer(&testDB{})
	// When
	last := "Last"
	req := &servicepb.AddNameRequest{GivenName: "First", Surname: &last}
	_, err := server.AddName(ctx, req)
	// Then
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

var (
	errUnimplemented = errors.New("unimplemented")
	errUnexpected    = errors.New("unexpected")
)

type testDB struct{}

func (*testDB) AddName(_ context.Context, name database.Name) error {
	if name.GivenName == "First" && *name.Surname == "Last" {
		return nil
	}

	return errUnexpected
}

func (*testDB) ListNames(context.Context, func(name database.Name) error) error {
	return errUnimplemented
}

func (*testDB) Disconnect(context.Context) (err error) {
	return errUnimplemented
}
