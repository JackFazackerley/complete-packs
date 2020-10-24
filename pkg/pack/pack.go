package pack

type Pack struct {
	Count float64 `json:"count"`
	Size  float64 `json:"size" db:"size"`
}
