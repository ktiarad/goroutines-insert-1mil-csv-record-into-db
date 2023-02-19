package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DbHost, DbPort, DbUser, DbPass, DbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(DbMaxConns)
	db.SetMaxIdleConns(DbMaxIdleConns)

	fmt.Println("Connection DB success")

	return db, nil
}
