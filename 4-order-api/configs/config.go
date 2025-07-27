package configs

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
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

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("некорректные данные в конфигурации: %v", err)
	}

	return &config, nil
}

func (c *Config) validate() error {
	var errs []error

	if c.Database.Host == "" {
		errs = append(errs, errors.New("поле host пустое"))
	}

	if c.Database.Port == "" {
		errs = append(errs, errors.New("поле host пустое"))
	}

	if c.Database.Password == "" {
		errs = append(errs, errors.New("поле password пустое"))
	}

	if c.Database.User == "" {
		errs = append(errs, errors.New("поле user пустое"))
	}

	if c.Database.Name == "" {
		errs = append(errs, errors.New("поле name пустое"))
	}

	if c.Database.SSLMode == "" {
		errs = append(errs, errors.New("поле ssl_mode пустое"))
	}

	return errors.Join(errs...)
}
