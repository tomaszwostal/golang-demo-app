package storage

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (c *Config) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Host, c.Port, c.User, c.DBName, c.Password, c.SSLMode)
}

var (
	gormOpen     = gorm.Open
	postgresOpen = postgres.Open
)

func NewConnection(config *Config) (*gorm.DB, error) {
	dsn := config.String()
	db, err := gormOpen(postgresOpen(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
