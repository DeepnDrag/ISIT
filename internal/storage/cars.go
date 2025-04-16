package storage

import (
	"ISIT/internal/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

// CarsRepo представляет репозиторий для работы с автомобилями.
type CarsRepo struct {
	db *gorm.DB
}

// NewCarsRepo создает новый экземпляр CarsRepo.
func NewCarsRepo(db *gorm.DB) *CarsRepo {
	return &CarsRepo{
		db: db,
	}
}

// Create добавляет новый автомобиль в базу данных и возвращает его ID.
func (r *CarsRepo) Create(car *models.Car) (uint, error) {
	if err := r.db.Create(car).Error; err != nil {
		return 0, err
	}
	log.Println("aaaaaaaaaaa", car)
	return car.ID, nil
}

// GetByID получает автомобиль по его ID.
func (r *CarsRepo) GetByID(id uint) (*models.Car, error) {
	var car models.Car

	// Выполняем запрос к базе данных для поиска автомобиля по ID
	if err := r.db.First(&car, id).Error; err != nil {
		// Если запись не найдена, возвращаем ошибку gorm.ErrRecordNotFound
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("car with id %d not found", id)
		}
		// Возвращаем другую ошибку, если произошла проблема с базой данных
		return nil, fmt.Errorf("failed to get car by id: %w", err)
	}

	return &car, nil
}

// GetAll получает список всех автомобилей.
func (r *CarsRepo) GetBrands() ([]string, error) {
	var brands []string

	// Используем DISTINCT для получения уникальных марок
	if err := r.db.Model(&models.Car{}).Distinct("brand").Pluck("make", &brands).Error; err != nil {
		return nil, err
	}

	return brands, nil
}

// Update обновляет данные автомобиля.
func (r *CarsRepo) Update(car *models.Car) error {
	return r.db.Save(car).Error
}

// Delete удаляет автомобиль по ID.
func (r *CarsRepo) Delete(id uint) error {
	return r.db.Delete(&models.Car{}, id).Error
}

func (r *CarsRepo) GetModelsByBrand(brand string) ([]string, error) {
	var model []string

	// Запрос для получения уникальных моделей для указанной марки
	if err := r.db.Model(&models.Car{}).Where("brand = ?", brand).Distinct("model").Pluck("model", &model).Error; err != nil {
		return nil, err
	}

	return model, nil
}

// FilterByStatus фильтрует автомобили по статусу (например, "available", "rented").
func (r *CarsRepo) FilterByStatus(status string) ([]models.Car, error) {
	var cars []models.Car
	if err := r.db.Where("status = ?", status).Preload("Location").Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

// FilterByLocation фильтрует автомобили по ID локации.
func (r *CarsRepo) FilterByLocation(locationID uint) ([]models.Car, error) {
	var cars []models.Car
	if err := r.db.Where("location_id = ?", locationID).Preload("Location").Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

func (r *CarsRepo) Filter(req models.SearchCarRequest) ([]models.Car, error) {
	var cars []models.Car

	// Начинаем запрос к таблице Car
	query := r.db.Model(&models.Car{}).
		Where("status = ?", "available") // Ищем только доступные автомобили

	// Добавляем фильтры из SearchCarRequest
	if req.Brand != "" {
		query = query.Where("brand = ?", req.Brand)
	}
	if req.Model != "" {
		query = query.Where("model = ?", req.Model)
	}
	if req.YearFrom > 0 {
		query = query.Where("year >= ?", req.YearFrom)
	}
	if req.YearTo > 0 {
		query = query.Where("year <= ?", req.YearTo)
	}
	if req.MinPrice > 0 {
		query = query.Where("price_per_day >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("price_per_day <= ?", req.MaxPrice)
	}

	// Проверяем доступность автомобиля на указанный период
	if req.StartDate != "" && req.EndDate != "" {
		query = query.Where("id NOT IN (?)", r.db.Model(&models.Order{}).
			Select("car_id").
			Where("status IN (?)", []string{"pending", "confirmed"}). // Учитываем активные заказы
			Where("(? BETWEEN start_date AND end_date OR ? BETWEEN start_date AND end_date)",
				req.StartDate, req.EndDate,
			),
		)
	}

	// Выполняем запрос
	result := query.Find(&cars)
	if result.Error != nil {
		return nil, result.Error
	}

	return cars, nil
}

// Search выполняет поиск автомобилей по ключевым словам (например, в модели или марке).
func (r *CarsRepo) Search(query string) ([]models.Car, error) {
	var cars []models.Car
	if err := r.db.Where("make LIKE ? OR model LIKE ?", "%"+query+"%", "%"+query+"%").
		Preload("Location").
		Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}

// GetAllWithPagination возвращает список автомобилей с пагинацией.
func (r *CarsRepo) GetAllWithPagination(page, limit int) ([]models.Car, int64, error) {
	var cars []models.Car
	var total int64

	offset := (page - 1) * limit
	query := r.db.Preload("Location").Offset(offset).Limit(limit)

	if err := query.Find(&cars).Error; err != nil {
		return nil, 0, err
	}
	if err := r.db.Model(&models.Car{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return cars, total, nil
}

// SortByPrice сортирует автомобили по цене (возрастание или убывание).
func (r *CarsRepo) SortByPrice(ascending bool) ([]models.Car, error) {
	var cars []models.Car
	order := "price_per_day ASC"
	if !ascending {
		order = "price_per_day DESC"
	}

	if err := r.db.Preload("Location").Order(order).Find(&cars).Error; err != nil {
		return nil, err
	}
	return cars, nil
}
