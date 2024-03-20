package service

import (
	"backend/pkg/model"
	"backend/pkg/repository"
	"context"
	"math/rand"
	"time"
)

type tasks struct {
	repo *repository.Repo
}

func NewTask(repo *repository.Repo) *tasks {
	return &tasks{
		repo: repo,
	}
}

func (s *tasks) GenerateUserVariance(ctx context.Context) (int, [][][]float64) {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(len(bankVariance)) + 1

	return randomNumber, bankVariance[randomNumber]
}

func (s *tasks) UpdateUserVariance(ctx context.Context, userId int, labId int, variance model.Variance) error {
	return s.repo.UpdateVariance(userId, labId, variance)
}

func (s *tasks) GetVariance(ctx context.Context, userId, labId int) (model.Variance, error) {
	return s.repo.GetVariance(userId, labId)
}
