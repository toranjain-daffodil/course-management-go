package repository

import "database/sql"

type studentRepositorySQL struct {
	db *sql.DB
}

func NewStudentRepositorySQL(db *sql.DB) *studentRepositorySQL {
	return &studentRepositorySQL{
		db: db,
	}
}