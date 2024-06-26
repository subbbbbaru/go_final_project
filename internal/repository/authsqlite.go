package repository

import (
	"database/sql"
)

type AuthSQLite struct {
	db *sql.DB
}

func NewAuthSQLite(db *sql.DB) *AuthSQLite {
	return &AuthSQLite{db: db}
}

func (r *AuthSQLite) CreateUser(user string) (int, error) {
	var id int
	// query := fmt.Sprintf("INSERT INTO %s (name, username, pwd_hash) values ($1, $2, $3) RETURNING id", usersTable)
	// row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	// if err := row.Scan(&id); err != nil {
	// 	return 0, err
	// }
	return id, nil
}

func (r *AuthSQLite) GetUser(username, password string) (string, error) {
	var user string
	var err error
	// query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND pwd_hash=$2", usersTable)
	// err := r.db.Get(&user, query, username, password)

	return user, err
}
