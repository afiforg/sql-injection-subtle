package repository

import (
	"database/sql"
	"sql-injection-subtle/internal/querybuilder"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername looks up users by username. The search term is passed
// to the query builder to form the WHERE clause—no concatenation here.
func (r *UserRepository) FindByUsername(username string) ([]User, error) {
	cond := querybuilder.BuildCondition("username", username)
	query := "SELECT id, username, email FROM users WHERE " + cond
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

// FindWithSort allows sorting by a column. Column and direction come
// from the service layer (ultimately from request params).
func (r *UserRepository) FindWithSort(sortColumn, sortDir string) ([]User, error) {
	orderClause := querybuilder.OrderBy(sortColumn, sortDir)
	query := "SELECT id, username, email FROM users " + orderClause
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUsers(rows)
}

type User struct {
	ID       int
	Username string
	Email    string
}

func scanUsers(rows *sql.Rows) ([]User, error) {
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, nil
}
