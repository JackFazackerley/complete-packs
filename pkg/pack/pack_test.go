package pack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPacks_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		p       Packs
		want    []byte
		wantErr error
	}{
		{
			name: "should remove object with count of 0",
			p: Packs{
				{
					Count: 1,
					Size:  10,
				},
				{
					Count: 0,
					Size:  20,
				},
			},
			want:    []byte(`[{"count":1,"size":10}]`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.MarshalJSON()
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
