package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"context"
)

type task interface {
	GenerateUserVariance1A(ctx context.Context) (int, [][][]float64)
	GenerateUserVariance1B(ctx context.Context) (int, [][][]float64)
	GenerateUserVariance2(ctx context.Context) (int, [][][]float64)
	UpdateUserVariance1A(ctx context.Context, userId int, variance model.Variance1A) error
	UpdateUserVariance1B(ctx context.Context, userId int, variance model.Variance1B) error
	UpdateUserVariance2(ctx context.Context, userId int, variance model.Variance2) error
	GetVariance1A(ctx context.Context, userId int) (model.Variance1A, error)
	GetVariance1B(ctx context.Context, userId int) (model.Variance1B, error)
	GetVariance2(ctx context.Context, userId int) (model.Variance2, error)
}

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
}

type lab1a interface {
	CheckLab1AImportanceMatrix(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error)
	CheckLab1AImportanceMatrixFirstCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error)
	CheckLab1AImportanceMatrixSecondCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error)
	CheckLab1AImportanceMatrixThirdCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error)
	CheckLab1AImportanceMatrixFourthCriteria(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error)
	CheckLab1AChosenAlternative(ctx context.Context, userId int, userMatrix [][]float64) (int, []float64, int, error)
}

type lab1b interface {
}

type lab2 interface {
}

type Service struct {
	external
	commonLab
	lab1a
	lab1b
	lab2
	task
}

func NewService(repo *repository.Repo) *Service {
	return &Service{
		external:  NewExternalService(),
		commonLab: NewCommonLabService(repo),
		lab1a:     NewLab1aService(repo),
		lab1b:     NewLab1bService(repo),
		lab2:      NewLab2Service(repo),
		task:      NewTask(repo),
	}
}
