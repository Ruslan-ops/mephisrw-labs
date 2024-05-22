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

const (
	lab1AId = 1
	lab1BId = 2
	lab2Id  = 3
)

func NewVarianceRepo(db *sqlx.DB) *VarianceRepo {
	return &VarianceRepo{db: db}
}

func (r *VarianceRepo) UpdateVariance(userId int, labId int, variance interface{}) error {
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

func (r *VarianceRepo) GetVariance(userId, labId int) (interface{}, error) {
	var data []byte
	query := fmt.Sprintf("SELECT variance FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&data, query, userId, labId); err != nil {
		return model.Variance1A{}, err
	}

	if labId == lab1AId {
		var variance model.Variance1A
		if err := json.Unmarshal(data, &variance); err != nil {
			return model.Variance1A{}, err
		}
		return variance, nil
	} else if labId == lab1BId {
		var variance model.Variance1B
		if err := json.Unmarshal(data, &variance); err != nil {
			return model.Variance1A{}, err
		}
		return variance, nil
	} else if labId == lab2Id {
		var variance model.Variance2
		if err := json.Unmarshal(data, &variance); err != nil {
			return model.Variance1A{}, err
		}
		return variance, nil
	}

	return nil, fmt.Errorf("variance not found for userId: %d, labId: %d", userId, labId)
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

func (r *VarianceRepo) CheckIsEmptyVariant(userId, labId int) bool {
	var data []byte
	query := fmt.Sprintf("SELECT variance FROM %s WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if err := r.db.Get(&data, query, userId, labId); err != nil {
		return true
	}

	if string(data) == "null" {
		return true
	}

	return false
}
