package config

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(dsn string) (*gorm.DB, *sql.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, nil, err
	}

	return gormDB, sqlDB, nil
}