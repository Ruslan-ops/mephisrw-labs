package repository

import (
	"backend/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type VarianceRepo struct {
	db *sqlx.DB
}

func NewVarianceRepo(db *sqlx.DB) *VarianceRepo {
	return &VarianceRepo{db: db}
}

func (r *VarianceRepo) UpdateVariance(userId int, labId int, variance model.Variance) error {
	jsonData, err := json.Marshal(variance)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET variance = $1 WHERE user_id = $2 AND internal_lab_id = $3", usersTable)
	if _, err = r.db.Exec(query, jsonData, userId, labId); err != nil {
		return err
	}

	return nil
}

func (r *VarianceRepo) GetVariance(userId, labId int) (model.Variance, error) {
	var data []byte
	query := fmt.Sprintf("SELECT variance FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&data, query, userId, labId); err != nil {
		return model.Variance{}, err
	}

	var variance model.Variance
	if err := json.Unmarshal(data, &variance); err != nil {
		return model.Variance{}, err
	}

	return variance, nil
}

func (r *VarianceRepo) CheckVariance(userId, labId int) error {
	var user int
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&user, query, userId, labId); err != nil {
		return err
	}

	if user == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
