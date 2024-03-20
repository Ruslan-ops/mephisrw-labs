package service

import (
	"context"
	"math/rand"
	"time"
)

type tasks struct{}

func NewTask() *tasks {
	return &tasks{}
}

func (s *tasks) GenerateUserVariance(ctx context.Context) (int, [][]float64) {
	rand.Seed(time.Now().UnixNano())

	randomNumber := rand.Intn(len(bankVariance)) + 1

	return randomNumber, bankVariance[randomNumber]
}