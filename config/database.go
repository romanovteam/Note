package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"notebook/logger"
	"notebook/repository"
)

// ConnectDatabase устанавливает соединение с базой данных и возвращает объект *gorm.DB
func ConnectDatabase(fileLogger logger.Logger) (*gorm.DB, error) {
	// Data Source Name
	dsn := "host=localhost user=postgres password=admin dbname=z port=5445 sslmode=disable TimeZone=Europe/Moscow"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fileLogger.LogError(err)
		return nil, err
	}
	fmt.Println("Database connected successfully.")

	// Миграция схемы
	err = db.AutoMigrate(&repository.Arg{})
	if err != nil {
		fileLogger.LogError(err)
		return nil, err
	}

	return db, nil
}
