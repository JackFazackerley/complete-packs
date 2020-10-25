package cache

import (
	"sync"
	"testing"

	"github.com/JackFazackerley/complete-packs/internal/database/sqlite"

	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/JackFazackerley/complete-packs/internal/database"
)

type TestDatbase struct {
	shouldError bool
	sizes       []float64
}

func (t *TestDatbase) ReadPacks() (packs []float64, err error) {
	if t.shouldError {
		return nil, sqlite.ReadError
	}
	return t.sizes, nil
}

func (t *TestDatbase) WritePack(size int) error {
	return nil
}

func (t *TestDatbase) DeletePack(size int) error {
	return nil
}

func (t *TestDatbase) Close() error {
	return nil
}

func TestSizes_Load(t *testing.T) {
	type fields struct {
		db        database.Database
		mux       *sync.RWMutex
		sizesASC  []float64
		sizesDESC []float64
	}
	type sizes struct {
		sizesASC  []float64
		sizesDESC []float64
	}
	tests := []struct {
		name          string
		fields        fields
		expectedErr   error
		expectedSizes sizes
	}{
		{
			name: "errors",
			fields: fields{
				db:        &TestDatbase{shouldError: true},
				mux:       &sync.RWMutex{},
				sizesASC:  []float64{},
				sizesDESC: []float64{},
			},
			expectedErr: sqlite.ReadError,
			expectedSizes: sizes{
				sizesASC:  []float64{},
				sizesDESC: []float64{},
			},
		},
		{
			name: "sorted sizes",
			fields: fields{
				db:        &TestDatbase{shouldError: false, sizes: []float64{10, 20, 30}},
				mux:       &sync.RWMutex{},
				sizesASC:  []float64{},
				sizesDESC: []float64{},
			},
			expectedErr: nil,
			expectedSizes: sizes{
				sizesASC:  []float64{10, 20, 30},
				sizesDESC: []float64{30, 20, 10, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Sizes{
				db:        tt.fields.db,
				mux:       tt.fields.mux,
				sizesASC:  tt.fields.sizesASC,
				sizesDESC: tt.fields.sizesDESC,
			}
			err := c.Load()
			assert.Equal(t, tt.expectedErr, errors.Cause(err))
			assert.Equal(t, tt.expectedSizes.sizesDESC, c.sizesDESC)
			assert.Equal(t, tt.expectedSizes.sizesASC, c.sizesASC)
		})
	}
}

func TestSizes_Asc(t *testing.T) {
	type fields struct {
		db        database.Database
		mux       *sync.RWMutex
		sizesASC  []float64
		sizesDESC []float64
	}
	tests := []struct {
		name     string
		fields   fields
		expected []float64
	}{
		{
			name: "asc returns correctly",
			fields: fields{
				db:        nil,
				mux:       &sync.RWMutex{},
				sizesASC:  []float64{10, 20, 30},
				sizesDESC: nil,
			},
			expected: []float64{10, 20, 30},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Sizes{
				db:        tt.fields.db,
				mux:       tt.fields.mux,
				sizesASC:  tt.fields.sizesASC,
				sizesDESC: tt.fields.sizesDESC,
			}
			got := c.Asc()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestSizes_Desc(t *testing.T) {
	type fields struct {
		db        database.Database
		mux       *sync.RWMutex
		sizesASC  []float64
		sizesDESC []float64
	}
	tests := []struct {
		name     string
		fields   fields
		expected []float64
	}{
		{
			name: "desc returns correctly",
			fields: fields{
				db:        nil,
				mux:       &sync.RWMutex{},
				sizesASC:  nil,
				sizesDESC: []float64{30, 20, 10, 0},
			},
			expected: []float64{30, 20, 10, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Sizes{
				db:        tt.fields.db,
				mux:       tt.fields.mux,
				sizesASC:  tt.fields.sizesASC,
				sizesDESC: tt.fields.sizesDESC,
			}
			got := c.Desc()
			assert.Equal(t, tt.expected, got)
		})
	}
}
