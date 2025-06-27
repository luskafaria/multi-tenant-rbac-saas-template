package database

import (
	"database/sql"
	_ "github.com/lib/pq"

	db "github.com/lucasfaria/rbac/api-go/internal/database/sqlc"
)

func NewDB(source string) (*sql.DB, *db.Queries, error) {
	dbConn, err := sql.Open("postgres", source)
	if err != nil {
		return nil, nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, nil, err
	}

	queries := db.New(dbConn)

	return dbConn, queries, nil
}
