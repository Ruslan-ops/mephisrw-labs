package model

type Variant1B struct {
	Number   int         `json:"number" db:"id"`
	Variance interface{} `json:"variance"`
}
