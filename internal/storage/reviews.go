package storage

import (
	"ISIT/internal/models"
	"gorm.io/gorm"
)

// ReviewsRepo представляет репозиторий для работы с отзывами.
type ReviewsRepo struct {
	db *gorm.DB
}

// NewReviewsRepo создает новый экземпляр ReviewsRepo.
func NewReviewsRepo(db *gorm.DB) *ReviewsRepo {
	return &ReviewsRepo{
		db: db,
	}
}

// Create добавляет новый отзыв в базу данных и возвращает его ID.
func (r *ReviewsRepo) Create(review *models.Review) (uint, error) {
	if err := r.db.Create(review).Error; err != nil {
		return 0, err
	}
	return review.ID, nil
}

// GetAll получает список всех отзывов.
func (r *ReviewsRepo) GetAll() ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Preload("User").Preload("Car").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// Update обновляет данные отзыва.
func (r *ReviewsRepo) Update(review *models.Review) error {
	return r.db.Save(review).Error
}

// Delete удаляет отзыв по его ID.
func (r *ReviewsRepo) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

// FilterByUserID фильтрует отзывы по ID пользователя.
func (r *ReviewsRepo) FilterByUserID(userID uint) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Where("user_id = ?", userID).Preload("User").Preload("Car").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}

// FilterByCarID фильтрует отзывы по ID автомобиля.
func (r *ReviewsRepo) FilterByCarID(carID uint) ([]models.Review, error) {
	var reviews []models.Review

	// Выполняем запрос к базе данных
	result := r.db.Where("car_id = ?", carID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error // Возвращаем ошибку, если запрос завершился неудачно
	}

	return reviews, nil
}

// FilterByRating фильтрует отзывы по оценке (например, только 5-звездочные).
func (r *ReviewsRepo) FilterByRating(rating int) ([]models.Review, error) {
	var reviews []models.Review
	if err := r.db.Where("rating = ?", rating).Preload("User").Preload("Car").Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
