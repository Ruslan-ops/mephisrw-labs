package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"context"
	"fmt"
)

const (
	Lab1AId = 1
	Lab1BId = 2
	Lab2Id  = 3
)

type commonLabService struct {
	repo *repository.Repo
}

func NewCommonLabService(repo *repository.Repo) *commonLabService {
	return &commonLabService{repo: repo}
}

func (s *commonLabService) GetLabResult(ctx context.Context, userId, labId int) (int, error) {
	mark, err := s.repo.GetCurrentMark(userId, labId)
	if err != nil {
		return 0, err
	}

	if err := s.repo.UpdateUserInfo(model.UserRepo{
		UserId:        userId,
		InternalLabId: labId,
		IsDone:        true,
		Percentage:    mark,
	}); err != nil {
		return mark, err
	}

	return mark, nil
}

func (s *commonLabService) GetUserVariance(ctx context.Context, userId, labId int) (interface{}, error) {
	return s.repo.GetVariance(userId, labId)
}

func (s *commonLabService) IncrementPercentageDone(ctx context.Context, userId, labId, mark int) error {
	return s.repo.IncrementMark(userId, labId, mark)
}

func (s *commonLabService) UpdateLabStep(ctx context.Context, userId, labId, step int) error {
	return s.repo.UpdateCurrentStep(userId, labId, step)
}

func (s *commonLabService) GetCurrentMark(userId, labId int) (int, error) {
	return s.repo.GetCurrentMark(userId, labId)
}

func (s *commonLabService) GetUserIdByToken(labId int, token string) (int, error) {
	id, err := s.repo.GetUserIdByToken(labId, token)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, fmt.Errorf("no user found with token")
	}

	return id, nil
}

func (s *commonLabService) SaveUserToken(userId, labId int, userHeader string) error {
	return s.repo.UpdateToken(userId, labId, userHeader)
}

func (s *commonLabService) GetUserInfo(userId, labId int) (model.UserRepo, error) {
	return s.repo.GetUserInfo(userId, labId)
}

func (s *commonLabService) OpenLabForStudent(ctx context.Context, userId, labId, externalLabId int) (bool, error) {
	val := model.UserRepo{
		UserId:        userId,
		InternalLabId: labId,
		ExternalLabId: externalLabId,
		IsDone:        false,
		Percentage:    0,
	}
	user, err := s.repo.GetUserInfo(userId, labId)
	if err != nil {
		if err := s.repo.InsertUserInfo(val); err != nil {
			return false, err
		}
		return false, nil
	} else {
		if err := s.repo.UpdateUserInfo(val); err != nil {
			return user.IsDone, err
		}
		if err := s.repo.UpdateVariance(userId, labId, nil); err != nil {
			return user.IsDone, err
		}
	}

	return user.IsDone, nil
}

func (s *commonLabService) ClearToken(userId, labId int) error {
	return s.repo.ClearToken(userId, labId)
}

func (s *commonLabService) GetLabCurrentStep(ctx context.Context, userId, labId int) (int, error) {
	return s.repo.GetCurrentStep(userId, labId)
}

func (s *commonLabService) CloseLabForStudent(ctx context.Context, userId, labId int) error {
	val := model.UserRepo{
		UserId:        userId,
		InternalLabId: labId,
		IsDone:        true,
		Percentage:    0,
	}
	user, err := s.repo.GetUserInfo(userId, labId)
	if err != nil || user == (model.UserRepo{}) {
		if err := s.repo.InsertUserInfo(val); err != nil {
			return err
		}
	} else {
		if err := s.repo.UpdateUserInfo(val); err != nil {
			return err
		}
	}

	return nil
}

func (s *commonLabService) IsEmptyToken(userId, labId int) bool {
	user, err := s.repo.GetUserInfo(userId, labId)
	if err != nil {
		return true
	}
	if user.Token == "" {
		return true
	}

	return false
}

func (s *commonLabService) UpdateUserVariance(userId, labId int, variance interface{}) error {
	return s.repo.UpdateVariance(userId, labId, variance)
}

func (s *commonLabService) CheckIsEmptyVariant(userId, labId int) bool {
	return s.repo.CheckIsEmptyVariant(userId, labId)
}
