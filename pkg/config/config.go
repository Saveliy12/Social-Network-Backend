package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Saveliy12/prod2/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	DB     Postgres
	Server Server
	log    logger.LoggerInterface
}

type Postgres struct {
	User     string
	Name     string
	SSLMode  string
	Password string
}

type Server struct {
	Port int
}

func New() (*Config, error) {
	cfg := new(Config)

	cfg.log = logger.GetLogger()

	err := godotenv.Load()

	if err != nil {
		cfg.log.Fatal("Error loading .env file")
	}

	// Загрузка значений из переменных окружения
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.SSLMode = os.Getenv("DB_SSLMODE")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	// Преобразование порта в int
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		return nil, errors.New("SERVER_PORT is empty")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SERVER_PORT: %w", err)
	}
	cfg.Server.Port = port

	return cfg, nil
}
