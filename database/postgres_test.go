package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fillmore-labs/name-service/database"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/ksqltest"
)

var (
	errProjectName = errors.New("expected \"test\" project name")
	errUnexpected  = errors.New("unexpected type")
)

func TestAddName(t *testing.T) {
	t.Parallel()
	// Given
	ctx := context.Background()
	mockDB := ksql.Mock{
		InsertFn: func(ctx context.Context, table ksql.Table, record interface{}) error {
			rec, err := ksqltest.StructToMap(record)
			if err != nil {
				return err
			}
			name, ok := rec["name"].(database.Name)
			if !ok {
				return errUnexpected
			} else if name.GivenName != "test" || name.Surname != nil {
				return errProjectName
			}

			return ksqltest.FillStructWith(record, map[string]interface{}{
				"id": 42,
			})
		},
	}
	db := database.NewKsqlDB(mockDB, nil)
	// When
	err := db.AddName(ctx, database.Name{GivenName: "test"})
	// Then
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
