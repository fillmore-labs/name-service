package database

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/vingarcia/ksql"
	kpgx "github.com/vingarcia/ksql/adapters/kpgx5"
)

const (
	tableName = "names"
	chunkSize = 10
)

//go:embed migrations/*.sql
var fs embed.FS

type ksqlDB struct {
	pg         ksql.Provider
	namesTable ksql.Table
	close      func() error
}

func NewKsqlDB(pg ksql.Provider, closeFn func() error) Database {
	return &ksqlDB{
		pg:         pg,
		namesTable: ksql.NewTable(tableName),
		close:      closeFn,
	}
}

func newPostgres(ctx context.Context, connectionString string) (db Database, err error) {
	// Opening, closing and re-opening of connections could be parallelized or avoided.
	err = ensurePGStructure(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	pg, err := kpgx.New(ctx, connectionString, ksql.Config{})
	if err != nil {
		return nil, err
	}

	return NewKsqlDB(pg, pg.Close), nil
}

func ensurePGStructure(_ context.Context, connectionString string) error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	defer func() {
		_ = sourceInstance.Close()
	}()

	db := pgx.Postgres{}
	databaseDrv, err := db.Open(connectionString)
	if err != nil {
		return err
	}
	defer func() {
		_ = databaseDrv.Close()
	}()

	m, _ := migrate.NewWithInstance("iofs", sourceInstance, "pgx5", databaseDrv)

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func (db *ksqlDB) Disconnect(_ context.Context) (err error) {
	return db.close()
}

type ksqlName struct {
	ID   int  `ksql:"id"`
	Name Name `ksql:"name,json"`
}

func (db *ksqlDB) AddName(ctx context.Context, name Name) (err error) {
	record := ksqlName{Name: name}
	if err := db.pg.Insert(ctx, db.namesTable, &record); err != nil {
		return fmt.Errorf("can't insert name: %w", err)
	}

	return nil
}

func (db *ksqlDB) ListNames(ctx context.Context, sendName func(name Name) error) error {
	return db.pg.QueryChunks(ctx, ksql.ChunkParser{
		Query:     "FROM " + tableName,
		Params:    []interface{}{},
		ChunkSize: chunkSize,
		ForEachChunk: func(names []ksqlName) error {
			for _, name := range names {
				if err := sendName(name.Name); err != nil {
					return fmt.Errorf("can't send name: %w", err)
				}
			}

			return nil
		},
	})
}
