package configs

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

type Config struct {
	Email    string
	Password string
	Address  string
}

func LoadConfig() (*Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("ошибка при загрузке конфигурации: %v", err)
	}

	config.Email = os.Getenv("EMAIL")
	config.Password = os.Getenv("PASSWORD")
	config.Address = os.Getenv("ADDRESS")

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("некорректные данные в конфигурации: %v", err)
	}

	return &config, nil
}

func (c *Config) validate() error {

	if c.Email == "" {
		return errors.New("поле email пустое")
	}

	match, err := regexp.MatchString(`[A-Za-z0-9\._%+\-]+@[A-Za-z0-9\.\-]+\.[A-Za-z]{2,}`, c.Email)
	if err != nil {
		return err
	}

	if !match {
		return errors.New("некорректный email")
	}

	if c.Password == "" {
		return errors.New("поле password пустое")
	}

	if c.Address == "" {
		return errors.New("поле address пустое")
	}

	return nil
}
