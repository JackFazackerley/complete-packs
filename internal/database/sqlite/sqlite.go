package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/pkg/errors"

	"github.com/JackFazackerley/complete-packs/internal/config"
	"github.com/JackFazackerley/complete-packs/internal/database"
)

const (
	driver = "sqlite3"

	GetQuery    = "SELECT size FROM packs"
	InsertQuery = "INSERT INTO packs (size) VALUES (?)"
	DeleteQuery = "DELETE FROM packs WHERE size = ?"
)

var (
	readError   = errors.New("reading packs")
	writeError  = errors.New("writing pack")
	deleteError = errors.New("deleting pack")
)

type SQLite struct {
	db *sqlx.DB
}

func New(config config.Database) (database.Database, error) {
	db, err := sqlx.Open(driver, config.Address())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to sqlite3 database")
	}

	return &SQLite{
		db: db,
	}, nil
}

// ReadPacks reads the packs into the passed slice of packs
func (s *SQLite) ReadPacks() (packs []float64, err error) {
	if err := s.db.Select(&packs, GetQuery); err != nil {
		return nil, errors.Wrap(readError, err.Error())
	}

	return packs, nil
}

func (s *SQLite) WritePack(size int) error {
	if _, err := s.db.Exec(InsertQuery, size); err != nil {
		return errors.Wrap(writeError, err.Error())
	}
	return nil
}

func (s *SQLite) DeletePack(size int) error {
	if _, err := s.db.Exec(DeleteQuery, size); err != nil {
		return errors.Wrap(deleteError, err.Error())
	}
	return nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
