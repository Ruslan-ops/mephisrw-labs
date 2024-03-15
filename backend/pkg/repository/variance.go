package repository

import (
	"github.com/jmoiron/sqlx"
)

type VarianceRepo struct {
	db *sqlx.DB
}

func NewVarianceRepo(db *sqlx.DB) *VarianceRepo {
	return &VarianceRepo{db: db}
}
