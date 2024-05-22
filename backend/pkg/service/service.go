package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"context"
)

type external interface {
	SendLabMark(ctx context.Context, userId, labId, percentage int) error
	GetUserId(ctx context.Context, token string) (int, error)
}

type commonLab interface {
	GetLabResult(ctx context.Context, userId, labId int) (int, error)
	IncrementPercentageDone(ctx context.Context, userId, labId, mark int) error
	UpdateLabStep(ctx context.Context, userId, labId, step int) error
	GetUserIdByToken(labId int, token string) (int, error)
	SaveUserToken(userId, labId int, userHeader string) error
	GetUserInfo(userId, labId int) (model.UserRepo, error)
	OpenLabForStudent(ctx context.Context, userId, labId, externalLabId int) (bool, error)
	CloseLabForStudent(ctx context.Context, userId, labId int) error
	ClearToken(userId, labId int) error
	GetLabCurrentStep(ctx context.Context, userId, labId int) (int, error)
	IsEmptyToken(userId, labId int) bool
	GetCurrentMark(userId, labId int) (int, error)
	UpdateUserVariance(userId, labId int, variance interface{}) error
	GetUserVariance(ctx context.Context, userId, labId int) (interface{}, error)
	CheckIsEmptyVariant(userId, labId int) bool
}

type lab1BVariance interface {
	GetIdealVariant1B() (model.Variant1B, error)
}

type Service struct {
	external
	commonLab
	lab1BVariance
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		external:      NewExternalService(),
		commonLab:     NewCommonLabService(repo),
		lab1BVariance: NewLab1BService(repo),
	}
}
