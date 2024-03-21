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
	UpdateVariance1A(userId int, labId int, variance model.Variance1A) error
	UpdateVariance1B(userId int, labId int, variance model.Variance1B) error
	UpdateVariance2(userId int, labId int, variance model.Variance2) error
	GetVariance1A(userId, labId int) (model.Variance1A, error)
	GetVariance1B(userId, labId int) (model.Variance1B, error)
	GetVariance2(userId, labId int) (model.Variance2, error)
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

type Repo struct {
	userRepo
	varianceRepo
	tokenRepo
	markRepo
}

func NewRepo(db *sqlx.DB) *Repo {
	return &Repo{
		userRepo:     NewUserRepo(db),
		tokenRepo:    NewTokenRepo(db),
		varianceRepo: NewVarianceRepo(db),
		markRepo:     NewMarkRepo(db),
	}
}
