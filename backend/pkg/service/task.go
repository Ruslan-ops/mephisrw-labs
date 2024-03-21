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

func (s *tasks) GenerateUserVariance1A(ctx context.Context) (int, [][][]float64) {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(len(bankVariance)) + 1

	return randomNumber, bankVariance[randomNumber]
}

func (s *tasks) GenerateUserVariance1B(ctx context.Context) (int, [][][]float64) {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(len(bankVariance)) + 1

	return randomNumber, bankVariance[randomNumber]
}

func (s *tasks) GenerateUserVariance2(ctx context.Context) (int, [][][]float64) {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(len(bankVariance)) + 1

	return randomNumber, bankVariance[randomNumber]
}

func (s *tasks) UpdateUserVariance1A(ctx context.Context, userId int, variance model.Variance1A) error {
	return s.repo.UpdateVariance1A(userId, Lab1AId, variance)
}

func (s *tasks) UpdateUserVariance1B(ctx context.Context, userId int, variance model.Variance1B) error {
	return s.repo.UpdateVariance1B(userId, Lab1BId, variance)
}

func (s *tasks) UpdateUserVariance2(ctx context.Context, userId int, variance model.Variance2) error {
	return s.repo.UpdateVariance2(userId, Lab2Id, variance)
}

func (s *tasks) GetVariance1A(ctx context.Context, userId int) (model.Variance1A, error) {
	return s.repo.GetVariance1A(userId, Lab1AId)
}

func (s *tasks) GetVariance1B(ctx context.Context, userId int) (model.Variance1B, error) {
	return s.repo.GetVariance1B(userId, Lab1BId)
}

func (s *tasks) GetVariance2(ctx context.Context, userId int) (model.Variance2, error) {
	return s.repo.GetVariance2(userId, Lab2Id)
}
