package repository

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cnf "main/pkg/config"
	model "main/pkg/model"
)

func PrepareDB(setings *cnf.Config) (*gorm.DB, error) {
	// Подключение к базе данных

	dsn := "host=" + setings.DBHost
	dsn += " user=" + setings.DBUser
	dsn += " password=" + setings.DBPasword
	dsn += " dbname=" + setings.DBName
	dsn += " port=" + setings.DBPort
	dsn += " sslmode=" + setings.DBSslmode
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Ошибка подключения к базе данных: %v", err)
		return nil, err
	}

	// Авто-миграция для создания таблицы если её не существует
	err = db.AutoMigrate(&model.Anomalies{})
	if err != nil {
		log.Printf("Ошибка миграции: %v", err)
		return nil, err
	}

	return db, nil
}
