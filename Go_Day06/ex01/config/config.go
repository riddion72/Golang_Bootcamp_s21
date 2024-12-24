package config

// package main

import (
	"encoding/json"
	// "fmt"
	"log"
	"os"
)

type Config struct {
	Host          string `json:"POSTGRES_HOST"`
	Port          string `json:"POSTGRES_PORT"`
	User          string `json:"POSTGRES_USER"`
	DBName        string `json:"POSTGRES_DBNAME"`
	SSLMode       string `json:"POSTGRES_SSLMODE"`
	Password      string `json:"POSTGRES_PASSWORD"`
	AdminName     string `json:"ADMIN_NAME"`
	AdminPassword string `json:"ADMIN_PASSWORD"`
}

const configFile string = "config/admin_credentials.txt"

func LoadConfig() (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		log.Println(err)
	}
	log.Println(config)
	return &config, nil
}
