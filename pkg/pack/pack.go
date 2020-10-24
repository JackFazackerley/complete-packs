package pack

import "encoding/json"

type Pack struct {
	Count float64 `json:"count"`
	Size  float64 `json:"size" db:"size"`
}

type Packs []Pack

func (p Packs) MarshalJSON() ([]byte, error) {
	var packs []Pack

	for _, pack := range p {
		if pack.Count > 0 {
			packs = append(packs, pack)
		}
	}

	return json.Marshal(packs)
}
