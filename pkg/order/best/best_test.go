package best

import (
	"testing"

	"github.com/JackFazackerley/complete-packs/pkg/pack"

	"github.com/stretchr/testify/assert"
)

func TestOrder_Calculate(t *testing.T) {
	tests := []struct {
		name     string
		sizes    []float64
		target   int
		expected []pack.Pack
	}{
		{
			name:   "1 should give 250",
			sizes:  []float64{250, 500, 1000, 2000, 5000},
			target: 1,
			expected: []pack.Pack{
				{
					Count: 1,
					Size:  250,
				},
				{
					Count: 0,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 0,
					Size:  2000,
				},
				{
					Count: 0,
					Size:  5000,
				},
			},
		},
		{
			name:   "250 should give 250",
			sizes:  []float64{250, 500, 1000, 2000, 5000},
			target: 250,
			expected: []pack.Pack{
				{
					Count: 1,
					Size:  250,
				},
				{
					Count: 0,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 0,
					Size:  2000,
				},
				{
					Count: 0,
					Size:  5000,
				},
			},
		},
		{
			name:   "251 should give 500",
			sizes:  []float64{250, 500, 1000, 2000, 5000},
			target: 251,
			expected: []pack.Pack{
				{
					Count: 0,
					Size:  250,
				},
				{
					Count: 1,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 0,
					Size:  2000,
				},
				{
					Count: 0,
					Size:  5000,
				},
			},
		},
		{
			name:   "501 should give 1x500 1x250",
			sizes:  []float64{250, 500, 1000, 2000, 5000},
			target: 501,
			expected: []pack.Pack{
				{
					Count: 1,
					Size:  250,
				},
				{
					Count: 1,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 0,
					Size:  2000,
				},
				{
					Count: 0,
					Size:  5000,
				},
			},
		},
		{
			name:   "12001 should give 2x5000 1x2000 1x250",
			sizes:  []float64{250, 500, 1000, 2000, 5000},
			target: 12001,
			expected: []pack.Pack{
				{
					Count: 1,
					Size:  250,
				},
				{
					Count: 0,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 1,
					Size:  2000,
				},
				{
					Count: 2,
					Size:  5000,
				},
			},
		},
		{
			name:   "1 should give 1x5",
			sizes:  []float64{5, 250, 500, 1000, 2000, 5000},
			target: 1,
			expected: []pack.Pack{
				{
					Count: 1,
					Size:  5,
				},
				{
					Count: 0,
					Size:  250,
				},
				{
					Count: 0,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 0,
					Size:  2000,
				},
				{
					Count: 0,
					Size:  5000,
				},
			},
		},
		{
			name:   "251 should give 1x250 1x5",
			sizes:  []float64{5, 250, 500, 1000, 2000, 5000},
			target: 251,
			expected: []pack.Pack{
				{
					Count: 1,
					Size:  5,
				},
				{
					Count: 1,
					Size:  250,
				},
				{
					Count: 0,
					Size:  500,
				},
				{
					Count: 0,
					Size:  1000,
				},
				{
					Count: 0,
					Size:  2000,
				},
				{
					Count: 0,
					Size:  5000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.target, tt.sizes)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func BenchmarkCalculate(b *testing.B) {
	sizes := []float64{250, 500, 1000, 2000, 5000}
	target := 999999999

	for n := 0; n < b.N; n++ {
		Calculate(target, sizes)
	}
}
