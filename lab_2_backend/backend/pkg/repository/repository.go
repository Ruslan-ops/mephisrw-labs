package repository

import (
	"backend/pkg/model"
	"github.com/jmoiron/sqlx"
)

type userRepo interface {
	GetUserInfo(userId, labId int) (model.UserRepo, error)
	InsertUserInfo(user model.UserRepo) error
	UpdateUserInfo(user model.UserRepo) error
}

type varianceRepo interface {
	UpdateVariance(userId int, labId int, variance interface{}) error
	GetVariance(userId, labId int) (interface{}, error)
	CheckIsEmptyVariant(userId, labId int) bool
}

type tokenRepo interface {
	UpdateToken(userId int, labId int, token string) error
	ClearToken(userId, labId int) error
	GetUserIdByToken(labId int, token string) (int, error)
}

type markRepo interface {
	UpdateCurrentStep(userId, labId, step int) error
	GetCurrentStep(userId, labId int) (int, error)
	IncrementMark(userId, labId, mark int) error
	GetCurrentMark(userId, labId int) (int, error)
}

type lab1BVariance interface {
	GetIdealVariant1B() (model.Variant1B, error)
}

type Repo struct {
	userRepo
	varianceRepo
	tokenRepo
	markRepo
	lab1BVariance
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		userRepo:      NewUserRepo(db),
		tokenRepo:     NewTokenRepo(db),
		varianceRepo:  NewVarianceRepo(db),
		markRepo:      NewMarkRepo(db),
		lab1BVariance: NewSecondLabRepo(db),
	}
}
