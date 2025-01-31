package repository

import (
	"backend/pkg/model"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SecondLabRepo struct {
	db *sqlx.DB
}

func NewSecondLabRepo(db *sqlx.DB) *SecondLabRepo {
	return &SecondLabRepo{
		db: db,
	}
}

func (r *SecondLabRepo) GetIdealVariant1B() (model.Variant1B, error) {
	data := struct {
		Number   int    `db:"id"`
		Variance []byte `db:"variance"`
	}{}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY RANDOM() LIMIT 1", bankVariance1B)
	if err := r.db.Get(&data, query); err != nil {
		return model.Variant1B{}, err
	}

	var variant model.Variant1B
	if err := json.Unmarshal(data.Variance, &variant.Variance); err != nil {
		return model.Variant1B{}, err
	}
	variant.Number = data.Number

	return variant, nil
}
