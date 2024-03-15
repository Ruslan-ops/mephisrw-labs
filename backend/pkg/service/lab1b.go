package service

import (
	"backend/pkg/repository"
)

const Lab1BId = 2

type lab1bService struct {
	repo *repository.Repo
}

func NewLab1bService(repo *repository.Repo) *lab1bService {
	return &lab1bService{
		repo: repo,
	}
}
