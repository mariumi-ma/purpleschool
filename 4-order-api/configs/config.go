package configs

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string `validate:"required"`
		Port     string `validate:"required"`
		User     string `validate:"required"`
		Password string `validate:"required"`
		Name     string `validate:"required"`
		SSLMode  string `validate:"required"`
	}
}

func LoadConfig() (*Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка при загрузке конфигурации: %v", err)
	}

	config.Database.Host = os.Getenv("DB_HOST")
	config.Database.Port = os.Getenv("DB_PORT")
	config.Database.User = os.Getenv("DB_USER")
	config.Database.Password = os.Getenv("DB_PASSWORD")
	config.Database.Name = os.Getenv("DB_NAME")
	config.Database.SSLMode = os.Getenv("DB_SSLMODE")

	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("некорректные данные в конфигурации: %v", err)
	}

	return &config, nil
}

func validate[T any](body *T) error {
	validate := validator.New()

	if err := validate.Struct(body); err != nil {
		return err
	}

	return nil
}
