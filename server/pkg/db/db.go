package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Database struct {
	Conn *sql.DB
}

func Connect() (*Database, error) {
	db, err := sql.Open("postgres", "postgres://root:root@localhost:5432/go-chat?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Database{Conn: db}, nil
}

func (d *Database) CloseConnection() {
	if err := d.Conn.Close(); err != nil {
		panic(err)
	}
}
