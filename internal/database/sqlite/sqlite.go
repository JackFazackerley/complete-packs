package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"

	"github.com/JackFazackerley/complete-packs/internal/database"
)

const (
	GetQuery    = "SELECT size FROM packs"
	InsertQuery = "INSERT INTO packs (size) VALUES (?)"
	DeleteQuery = "DELETE FROM packs WHERE size = ?"
)

var (
	ReadError   = errors.New("reading packs")
	WriteError  = errors.New("writing pack")
	DeleteError = errors.New("deleting pack")
)

type SQLite struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) database.Database {
	return &SQLite{
		db: db,
	}
}

// ReadPacks reads the packs into the passed slice of packs
func (s *SQLite) ReadPacks() (packs []float64, err error) {
	if err := s.db.Select(&packs, GetQuery); err != nil {
		return nil, errors.Wrap(ReadError, err.Error())
	}

	return packs, nil
}

func (s *SQLite) WritePack(size int) error {
	if _, err := s.db.Exec(InsertQuery, size); err != nil {
		return errors.Wrap(WriteError, err.Error())
	}
	return nil
}

func (s *SQLite) DeletePack(size int) error {
	if _, err := s.db.Exec(DeleteQuery, size); err != nil {
		return errors.Wrap(DeleteError, err.Error())
	}
	return nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
