package database

import (
	"ISIT/internal/config"
	"ISIT/internal/models"
	"ISIT/internal/utils"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func Connection(config config.DB) (*gorm.DB, error) {
	if config.User == "" || config.Password == "" || config.Host == "" || config.Port == "" || config.Name == "" {
		return nil, fmt.Errorf("invalid database configuration")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Host, config.User, config.Password, config.Name, config.Port)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	var db *gorm.DB
	var err error
	for attempts := 0; attempts < 3; attempts++ {
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("gorm open error after retries: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}

	return db, nil
}

func Migrations(db *gorm.DB) error {

	// Создание таблиц
	err := db.AutoMigrate(&models.User{}, &models.Car{}, &models.Order{}, &models.Location{}, &models.Review{})
	if err != nil {
		return fmt.Errorf("ошибка при миграции: %w", err)
	}

	// Проверяем, есть ли уже администратор
	var existingAdmin models.User
	result := db.Where("role = ?", "admin").First(&existingAdmin)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		ph, _ := utils.HashPassword("admin123")
		// Если админа нет, создаем его
		adminUser := models.User{
			FirstName:    "Admin",
			LastName:     "User",
			Email:        "admin",
			PasswordHash: ph,
			Role:         "admin",
			CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		}

		// Добавляем в базу
		if err := db.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("ошибка при создании администратора: %w", err)
		}

		fmt.Println("Администратор создан: admin@example.com / admin123")
	}

	return nil
}
