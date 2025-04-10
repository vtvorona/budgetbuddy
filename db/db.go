package db

import (
	"budgetbuddy/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Переменная для базы данных
var DB *gorm.DB

// Инициализация базы данных
func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Автоматическое создание таблицы
	err = DB.AutoMigrate(&models.User{}, &models.Expense{}, models.Category{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

// Закрытие базы данных
func CloseDatabase() {
	sqlDatabase, err := DB.DB() // Получаем доступ к базовым объектам SQL
	if err != nil {
		log.Fatal("Failed to get database:", err)
	}
	if err := sqlDatabase.Close(); err != nil {
		log.Fatal("Failed to close database:", err)
	}
}
