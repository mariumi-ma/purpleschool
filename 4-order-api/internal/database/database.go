package database

import (
	"fmt"

	"purpleschool/configs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Db struct {
	*gorm.DB
}

func NewDB(config *configs.Config) (*Db, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Port,
		config.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Db{db}, nil
}
