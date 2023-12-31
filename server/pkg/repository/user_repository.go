package repository

import (
	"context"
	"database/sql"
	"server/pkg/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := "INSERT INTO users(username, email, password) VALUES($1, $2, $3) returning id"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	if err = stmt.QueryRow(user.Username, user.Email, user.Password).Scan(&user.Id); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := "SELECT * FROM users WHERE email = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var user model.User

	if err := stmt.QueryRow(email).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}
