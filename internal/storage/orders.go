package storage

import (
	"ISIT/internal/models"
	"errors"
	"gorm.io/gorm"
)

// OrdersRepo представляет репозиторий для работы с заказами.
type OrdersRepo struct {
	db *gorm.DB
}

// NewOrdersRepo создает новый экземпляр OrdersRepo.
func NewOrdersRepo(db *gorm.DB) *OrdersRepo {
	return &OrdersRepo{
		db: db,
	}
}

// Create добавляет новый заказ в базу данных и возвращает его ID.
func (r *OrdersRepo) Create(order *models.Order) (uint, error) {
	if err := r.db.Create(order).Error; err != nil {
		return 0, err
	}
	return order.ID, nil
}

// GetByUserID получает все заказы по ID пользователя.
func (r *OrdersRepo) GetByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order

	if err := r.db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Order{}, nil
		}
		return nil, err
	}

	return orders, nil
}

// GetAll получает список всех заказов.
func (r *OrdersRepo) GetAll() ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Preload("User").Preload("Car").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// Update обновляет данные заказа.
func (r *OrdersRepo) Update(order *models.Order) error {
	return r.db.Save(order).Error
}

// Delete удаляет заказ по его ID.
func (r *OrdersRepo) Delete(id uint) error {
	return r.db.Delete(&models.Order{}, id).Error
}

// FilterByStatus фильтрует заказы по статусу (например, "pending", "confirmed").
func (r *OrdersRepo) FilterByStatus(status string) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("status = ?", status).Preload("User").Preload("Car").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// FilterByUserID фильтрует заказы по ID пользователя.
func (r *OrdersRepo) FilterByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("user_id = ?", userID).Preload("User").Preload("Car").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// FilterByCarID фильтрует заказы по ID автомобиля.
func (r *OrdersRepo) FilterByCarID(carID uint) ([]models.Order, error) {
	var orders []models.Order
	if err := r.db.Where("car_id = ?", carID).Preload("User").Preload("Car").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
