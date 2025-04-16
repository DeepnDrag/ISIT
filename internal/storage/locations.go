package storage

import (
	"ISIT/internal/models"
	"gorm.io/gorm"
)

// LocationsRepo представляет репозиторий для работы с локациями.
type LocationsRepo struct {
	db *gorm.DB
}

// NewLocationsRepo создает новый экземпляр LocationsRepo.
func NewLocationsRepo(db *gorm.DB) *LocationsRepo {
	return &LocationsRepo{
		db: db,
	}
}

// Create добавляет новую локацию в базу данных и возвращает её ID.
func (r *LocationsRepo) Create(location *models.Location) (uint, error) {
	if err := r.db.Create(location).Error; err != nil {
		return 0, err
	}
	return location.ID, nil
}

// GetByID получает локацию по её ID.
func (r *LocationsRepo) GetByID(id uint) (*models.Location, error) {
	var location models.Location
	if err := r.db.First(&location, id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

// GetAll получает список всех локаций.
func (r *LocationsRepo) GetAll() ([]models.Location, error) {
	var locations []models.Location
	if err := r.db.Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

// Update обновляет данные локации.
func (r *LocationsRepo) Update(location *models.Location) error {
	return r.db.Save(location).Error
}

// Delete удаляет локацию по её ID.
func (r *LocationsRepo) Delete(id uint) error {
	return r.db.Delete(&models.Location{}, id).Error
}

// FilterByCity фильтрует локации по названию города.
func (r *LocationsRepo) FilterByCity(city string) ([]models.Location, error) {
	var locations []models.Location
	if err := r.db.Where("city = ?", city).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

// FilterByCountry фильтрует локации по названию страны.
func (r *LocationsRepo) FilterByCountry(country string) ([]models.Location, error) {
	var locations []models.Location
	if err := r.db.Where("country = ?", country).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}
