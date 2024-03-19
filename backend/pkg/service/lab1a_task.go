package service 

type Lab1ATask  struct {
	Number int64 `json:"number"`
	Task [][]float64 `json:"task"`
}

var Lab1ABankVariance = [][][]float64 {
	[][]float64 {
		[]float64 {1, 0.5, 0.5},
		[]float64 {0.5, 1, 0.5},
		[]float64 {0.5, 0.5, 1},
	},
}