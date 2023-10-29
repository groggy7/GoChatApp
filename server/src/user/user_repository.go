package user

import (
	"context"
	"database/sql"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *User) (*User, error) {
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

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := "SELECT * FROM users WHERE email = $1"

	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var user User

	if err := stmt.QueryRow(email).Scan(&user.Id, &user.Username, &user.Email, &user.Password); err != nil {
		return nil, err
	}

	return &user, nil
}
