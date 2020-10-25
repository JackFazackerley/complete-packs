package sqlite

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/jmoiron/sqlx"
)

const (
	tableDDL = `
	CREATE TABLE "packs" (
   		size INTEGER NOT NULL UNIQUE
	);`

	insert = `INSERT INTO packs ("size") VALUES (250),(500),(1000),(2000),(5000);`
)

func dbInit(create bool, t *testing.T) *sqlx.DB {
	t.Helper()

	db, err := sqlx.Open("sqlite3", ":memory:")
	if err != nil {
		t.Errorf("creating db: %s", err)
		return nil
	}

	if create {
		if _, err := db.Exec(tableDDL); err != nil {
			t.Errorf("creating table: %s", err)
			return nil
		}

		if _, err := db.Exec(insert); err != nil {
			t.Errorf("inserting into table: %s", err)
			return nil
		}
	}

	return db
}

func TestSQLite_ReadPacks(t *testing.T) {
	tests := []struct {
		name          string
		db            *sqlx.DB
		expectedPacks []float64
		expectedErr   error
	}{
		{
			name:          "returns packs",
			db:            dbInit(true, t),
			expectedPacks: []float64{250, 500, 1000, 2000, 5000},
			expectedErr:   nil,
		},
		{
			name:          "errors",
			db:            dbInit(false, t),
			expectedPacks: nil,
			expectedErr:   ReadError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQLite{
				db: tt.db,
			}
			defer s.Close()

			got, err := s.ReadPacks()
			assert.Equal(t, tt.expectedErr, errors.Cause(err))
			assert.Equal(t, tt.expectedPacks, got)
		})
	}
}

func TestSQLite_WritePack(t *testing.T) {
	tests := []struct {
		name        string
		db          *sqlx.DB
		size        int
		expectedErr error
	}{
		{
			name:        "should insert",
			db:          dbInit(true, t),
			size:        10,
			expectedErr: nil,
		},
		{
			name:        "already exists",
			db:          dbInit(true, t),
			size:        250,
			expectedErr: DuplicateKey,
		},
		{
			name:        "unknown error",
			db:          dbInit(false, t),
			size:        250,
			expectedErr: WriteError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQLite{
				db: tt.db,
			}
			defer s.Close()

			err := s.WritePack(tt.size)
			assert.Equal(t, tt.expectedErr, errors.Cause(err))
		})
	}
}

func TestSQLite_DeletePack(t *testing.T) {
	tests := []struct {
		name        string
		db          *sqlx.DB
		size        int
		expectedErr error
	}{
		{
			name:        "delete size",
			db:          dbInit(true, t),
			size:        250,
			expectedErr: nil,
		},
		{
			name:        "shouldn't error, size doesn't exist",
			db:          dbInit(true, t),
			size:        10,
			expectedErr: nil,
		},
		{
			name:        "database doesn't exist",
			db:          dbInit(false, t),
			size:        250,
			expectedErr: DeleteError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SQLite{
				db: tt.db,
			}
			err := s.DeletePack(tt.size)
			assert.Equal(t, tt.expectedErr, errors.Cause(err))
		})
	}
}
