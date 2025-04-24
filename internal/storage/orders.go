package storage

import (
	"ISIT/internal/models"
	"gorm.io/gorm"
	"log"
	"log/slog"
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
	slog.Info("заказ создан")
	return order.ID, nil
}

func (r *OrdersRepo) GetByUserID(userID uint) ([]models.OrderWithCarRequest, error) {
	var orders []models.OrderWithCarRequest

	// Выполняем запрос с JOIN для получения данных о машине
	if err := r.db.Table("\"order\"").
		Select("\"order\".id, \"order\".user_id, \"order\".car_id, car.brand as car_brand, car.model as car_model, \"order\".start_date, \"order\".end_date, \"order\".total_cost, \"order\".status, \"order\".created_at, \"order\".updated_at").
		Joins("JOIN car ON \"order\".car_id = car.id").
		Where("\"order\".user_id = ?", userID).
		Scan(&orders).Error; err != nil {

		// Логируем ошибку
		log.Println("Ошибка при получении заказов:", err)

		// Возвращаем ошибку в случае проблем
		return nil, err
	}

	// Если заказов нет, возвращаем пустой массив
	if len(orders) == 0 {
		return []models.OrderWithCarRequest{}, nil
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
