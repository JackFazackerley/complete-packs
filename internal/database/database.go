package database

import (
	"io"
)

type Database interface {
	ReadPacks() (packs []float64, err error)
	WritePack(size int) error
	DeletePack(size int) error

	io.Closer
}
