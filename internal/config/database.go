package config

import (
	"database/sql"
	"fmt"
	"time"
)

const DB_HOST_KEY = "DB_HOST"
const DB_PORT_KEY = "DB_PORT"
const DB_PASSWORD_KEY = "DB_PASSWORD"
const DB_USER_KEY = "DB_USER"
const DB_NAME_KEY = "DB_NAME"

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDBConfig() *DBConfig {
	dbHost := getEnv(DB_HOST_KEY, "localhost")
	dbPort := getEnv(DB_PORT_KEY, "3306")
	dbUser := getEnv(DB_USER_KEY, "root")
	dbPassword := getEnv(DB_PASSWORD_KEY, "root")
	dbName := getEnv(DB_NAME_KEY, "pubplay")

	return &DBConfig{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
	}
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func NewDB(config *DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.DSN())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
