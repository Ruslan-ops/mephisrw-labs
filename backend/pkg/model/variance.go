package model

type Variance struct {
	Number int           `json:"number"`
	Data   [][][]float64 `json:"data"`
}
