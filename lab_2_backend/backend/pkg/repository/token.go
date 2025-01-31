package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
)

type TokenRepo struct {
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (r *TokenRepo) UpdateToken(userId int, labId int, token string) error {
	query := fmt.Sprintf("UPDATE %s SET token = $1 WHERE user_id = $2 AND internal_lab_id = $3", usersTable)
	if _, err := r.db.Exec(query, token, userId, labId); err != nil {
		return err
	}

	return nil
}

type Claims struct {
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
	UserId int   `json:"user_id"`
}

func (c *Claims) Valid() error {
	if c.Exp < time.Now().Unix() {
		return errors.New("token is expired")
	}
	return nil
}

func (r *TokenRepo) GetUserIdByToken(labId int, token string) (int, error) {
	// Parse the token
	token = strings.TrimPrefix(token, "Bearer ")
	parser := jwt.Parser{}
	claims := &Claims{}
	_, _, err := parser.ParseUnverified(token, claims)
	if err != nil {
		return 0, errors.New("invalid token format")
	}

	// Validate claims (optional, since we're not fully parsing)
	if err := claims.Valid(); err != nil {
		return 0, err
	}

	// Return the user ID
	return claims.UserId, nil
}

func (r *TokenRepo) ClearToken(userId, labId int) error {
	query := fmt.Sprintf("UPDATE %s SET token = '' WHERE user_id = $1 AND internal_lab_id = $2", usersTable)
	if _, err := r.db.Exec(query, userId, labId); err != nil {
		return err
	}

	return nil
}
