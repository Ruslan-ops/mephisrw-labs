package service

import (
	"backend/pkg/repository"
)

const Lab1AId = 1

type lab1aService struct {
	repo *repository.Repo
}

func NewLab1aService(repo *repository.Repo) *lab1aService {
	return &lab1aService{repo: repo}
}
