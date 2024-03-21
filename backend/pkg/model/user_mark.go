package model

type UserLabMark struct {
	UserId     int `json:"user_id"`
	LabId      int `json:"laboratory_id"`
	Percentage int `json:"percentage"`
}

type UserLabMarkToResponse struct {
	UserId     int        `json:"user_id"`
	LabId      int        `json:"laboratory_id"`
	Variance   Variance1A `json:"variance"`
	Step       int        `json:"step"`
	Percentage int        `json:"percentage"`
}

type UserRepo struct {
	UserId        int    `json:"user_id" db:"user_id"`
	InternalLabId int    `json:"internal_lab_id" db:"internal_lab_id"`
	ExternalLabId int    `json:"external_lab_id" db:"external_lab_id"`
	IsDone        bool   `json:"is_done" db:"is_done"`
	Percentage    int    `json:"percentage" db:"percentage"`
	Token         string `json:"token" db:"token"`
}

type UserVarianceRepo struct {
	UserId   int        `json:"user_id" db:"user_id"`
	Variance Variance1A `json:"variance"`
}

type UserLabKey struct {
	UserId int `json:"user_id"`
	LabId  int `json:"lab_id"`
}

type UserIsDone struct {
	IsDone   bool       `json:"is_done"`
	Variance Variance1A `json:"variance"`
}
