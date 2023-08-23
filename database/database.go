package database

import (
	"context"
	"fmt"
	"os"
)

type Name struct {
	GivenName string  `json:"givenName"`
	Surname   *string `json:"surname"`
}

type Database interface {
	Disconnect(ctx context.Context) error
	AddName(ctx context.Context, name Name) error
	ListNames(ctx context.Context, sendName func(name Name) error) error
}

func NewDatabase(ctx context.Context) (db Database, err error) {
	// postgresql://username:password@host:port/database_name
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL != "" {
		fmt.Println("Using PostgreSQL")

		return newPostgres(ctx, postgresURL)
	}

	fmt.Println("WARNING: Using In-Memory Database")

	return newMem(ctx)
}
