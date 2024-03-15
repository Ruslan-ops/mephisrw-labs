package service

import (
	"backend/pkg/repository"
)

const Lab2Id = 3

type lab2Service struct {
	repo *repository.Repo
}

func NewLab2Service(repo *repository.Repo) *lab2Service {
	return &lab2Service{
		repo: repo,
	}
}
