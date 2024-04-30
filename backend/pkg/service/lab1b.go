package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
)

type lab1BService struct {
	repo *repository.Repo
}

func NewLab1BService(repo *repository.Repo) *lab1BService {
	return &lab1BService{
		repo: repo,
	}
}

func (s *lab1BService) GetIdealVariant1B() (model.Variant1B, error) {
	return s.repo.GetIdealVariant1B()
}
