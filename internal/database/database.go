package database

import (
	"io"

	"github.com/JackFazackerley/complete-packs/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	driver = "sqlite3"
)

type Database interface {
	ReadPacks() (packs []float64, err error)
	WritePack(size int) error
	DeletePack(size int) error

	io.Closer
}

func OpenDatabase(config config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, config.Address())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to sqlite3 database")
	}

	return db, nil
}
