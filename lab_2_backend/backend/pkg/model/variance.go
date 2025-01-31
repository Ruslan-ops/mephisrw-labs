package model

type Variance1A struct {
	Number int         `json:"number" binding:"required"`
	Data   interface{} `json:"data" binding:"required"`
}

type Variance1B struct {
	Number int         `json:"number" binding:"required"`
	Data   interface{} `json:"data" binding:"required"`
}

type Variance2 struct {
	Number int         `json:"number" binding:"required"`
	Data   interface{} `json:"data" binding:"required"`
}
