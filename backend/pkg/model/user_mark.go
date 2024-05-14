package model

type UserLabMark struct {
	UserId     int `json:"user_id"`
	LabId      int `json:"laboratory_id"`
	Percentage int `json:"percentage"`
}

type UserRepo struct {
	UserId        int    `json:"user_id" db:"user_id"`
	InternalLabId int    `json:"internal_lab_id" db:"internal_lab_id"`
	ExternalLabId int    `json:"external_lab_id" db:"external_lab_id"`
	IsDone        bool   `json:"is_done" db:"is_done"`
	Percentage    int    `json:"percentage" db:"percentage"`
	Token         string `json:"token" db:"token"`
}

type UserStepPercentage struct {
	Step       int `json:"step" binding:"required"`
	Percentage int `json:"percentage"`
}

type SendUserResult struct {
	Percentage int `json:"percentage"`
}
