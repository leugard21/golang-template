package db

import (
	"database/sql"
	"fmt"
	"log"
	"template/configs"

	_ "github.com/lib/pq"
)

func NewPostgresStorage(cfg configs.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("DB: Successfully connected!")
	return db, nil
}
