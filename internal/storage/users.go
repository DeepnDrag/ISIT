package storage

import (
	"ISIT/internal/models"
	"fmt"
	"gorm.io/gorm"
	"time"
)

// UsersRepo представляет репозиторий для работы с пользователями.
type UsersRepo struct {
	db *gorm.DB
}

// NewUsersRepo создает новый экземпляр UsersRepo.
func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

// Create создает нового пользователя в базе данных на основе email и password.
func (r *UsersRepo) Create(email, password string) (*models.User, error) {
	// Создаем нового пользователя с минимальными данными
	user := models.User{
		Email:        email,
		PasswordHash: password,   // Предполагается, что пароль уже хэширован перед вызовом метода
		Role:         "customer", // Устанавливаем роль по умолчанию
		CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:    time.Now().Format("2006-01-02 15:04:05"),
	}

	// Добавляем пользователя в базу данных
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail получает пользователя по его email.
func (r *UsersRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Специальная обработка для случая, когда пользователь не найден
			return nil, fmt.Errorf("пользователь с email %s не найден", email)
		}
		// Другие ошибки от базы данных
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}
	return &user, nil
}

func (r *UsersRepo) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAll получает список всех пользователей.
func (r *UsersRepo) GetAll() ([]models.User, error) {
	var users []models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// Update обновляет данные пользователя.
func (r *UsersRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete удаляет пользователя по ID.
func (r *UsersRepo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}
