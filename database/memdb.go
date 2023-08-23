package database

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	memdb "github.com/hashicorp/go-memdb"
)

type memDB struct {
	md      *memdb.MemDB
	counter atomic.Int32
}

const (
	namesTable = "names"
)

var errInvalidType = errors.New("name service: invalid type")

type memName struct {
	ID int32
	Name
}

func newMem(context.Context) (db Database, err error) {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			namesTable: {
				Name: namesTable,
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "ID"},
					},
				},
			},
		},
	}

	md, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	result := &memDB{md: md}
	result.counter.Store(1)

	return result, nil
}

func (*memDB) Disconnect(context.Context) (err error) {
	return nil
}

func (db *memDB) AddName(_ context.Context, name Name) (err error) {
	c := db.counter.Add(1)

	t := db.md.Txn(true)

	if err = t.Insert(namesTable, memName{ID: c, Name: name}); err != nil {
		t.Abort()

		return fmt.Errorf("can't add project: %w", err)
	}

	t.Commit()

	return nil
}

func (db *memDB) ListNames(_ context.Context, sendName func(name Name) error) error {
	t := db.md.Txn(false)

	iter, err := t.Get(namesTable, "id")
	if err != nil {
		t.Abort()

		return fmt.Errorf("can't set names: %w", err)
	}

	for {
		next := iter.Next()
		if next == nil {
			break
		}

		name, ok := next.(memName)
		if !ok {
			return errInvalidType
		}

		if err := sendName(name.Name); err != nil {
			return fmt.Errorf("can't send name: %w", err)
		}
	}

	t.Commit()

	return nil
}
