package config

import (
	"log"
	"os"

	env "github.com/joho/godotenv"
)

type Config struct {
	DBHost    string
	DBPort    string
	DBName    string
	DBUser    string
	DBPasword string
	DBSslmode string

	ServerHost string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	err := env.Load()
	if err != nil {
		log.Println("Error loading .env file")
		return &Config{}, err
	}

	fromEnv := &Config{
		DBHost:    os.Getenv("DB_HOST"),
		DBPort:    os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_NAME"),
		DBUser:    os.Getenv("DB_USER"),
		DBPasword: os.Getenv("DB_PASSWORD"),
		DBSslmode: os.Getenv("DB_SSLMODE"),

		ServerHost: os.Getenv("SERVER_TRANSMITTER_HOST"),
		ServerPort: os.Getenv("SERVER_TRANSMITTER_PORT"),
	}

	return fromEnv, nil
}
