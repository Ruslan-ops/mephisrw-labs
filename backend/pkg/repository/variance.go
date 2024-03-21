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

func (r *VarianceRepo) UpdateVariance1A(userId int, labId int, variance model.Variance1A) error {
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

func (r *VarianceRepo) UpdateVariance1B(userId int, labId int, variance model.Variance1B) error {
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

func (r *VarianceRepo) UpdateVariance2(userId int, labId int, variance model.Variance2) error {
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

func (r *VarianceRepo) GetVariance1A(userId, labId int) (model.Variance1A, error) {
	var data []byte
	query := fmt.Sprintf("SELECT variance FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&data, query, userId, labId); err != nil {
		return model.Variance1A{}, err
	}

	var variance model.Variance1A
	if err := json.Unmarshal(data, &variance); err != nil {
		return model.Variance1A{}, err
	}

	return variance, nil
}

func (r *VarianceRepo) GetVariance1B(userId, labId int) (model.Variance1B, error) {
	var data []byte
	query := fmt.Sprintf("SELECT variance FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&data, query, userId, labId); err != nil {
		return model.Variance1B{}, err
	}

	var variance model.Variance1B
	if err := json.Unmarshal(data, &variance); err != nil {
		return model.Variance1B{}, err
	}

	return variance, nil
}

func (r *VarianceRepo) GetVariance2(userId, labId int) (model.Variance2, error) {
	var data []byte
	query := fmt.Sprintf("SELECT variance FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&data, query, userId, labId); err != nil {
		return model.Variance2{}, err
	}

	var variance model.Variance2
	if err := json.Unmarshal(data, &variance); err != nil {
		return model.Variance2{}, err
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
