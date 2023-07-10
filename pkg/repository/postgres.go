package repository

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewConnection(config *Config) (*gorm.DB, error) {
	dsm := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsm), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
