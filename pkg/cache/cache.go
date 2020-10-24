package cache

import (
	"sort"
	"sync"

	"github.com/pkg/errors"

	"github.com/JackFazackerley/complete-packs/internal/database"
)

var (
	readError = errors.New("reading packs")
)

type Cache interface {
	Load() error
	Asc() []float64
	Desc() []float64
}

// Sizes is used to hold the pack sizes for both order versions
type Sizes struct {
	db database.Database

	mux *sync.RWMutex

	sizesASC  []float64
	sizesDESC []float64
}

func New(db database.Database) Cache {
	return &Sizes{
		db:  db,
		mux: &sync.RWMutex{},
	}
}

// Load quite literally loads the cache
func (c *Sizes) Load() error {
	packs, err := c.db.ReadPacks()
	if err != nil {
		return errors.Wrap(readError, err.Error())
	}

	// we want to block to stop any new calls possibly getting an empty slice
	c.mux.Lock()

	c.sizesASC = make([]float64, len(packs))
	c.sizesDESC = make([]float64, len(packs))

	copy(c.sizesASC, packs)
	copy(c.sizesDESC, packs)

	// sort for best
	sort.Slice(c.sizesASC, func(i, j int) bool {
		return c.sizesASC[i] < c.sizesASC[j]
	})

	// sort for fast
	c.sizesDESC = append(c.sizesDESC, 0.)
	sort.Slice(c.sizesDESC, func(i, j int) bool {
		return c.sizesDESC[i] > c.sizesDESC[j]
	})

	// unlock once we're done writing
	c.mux.Unlock()

	return nil
}

// Asc return the ascending pack sizes
func (c *Sizes) Asc() []float64 {
	c.mux.RLock()

	defer c.mux.RUnlock()

	return c.sizesASC
}

// Desc return the descending pack sizes
func (c *Sizes) Desc() []float64 {
	c.mux.RLock()

	defer c.mux.RUnlock()

	return c.sizesDESC
}
