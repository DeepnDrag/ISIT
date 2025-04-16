package storage

import (
	"ISIT/internal/models"
	"gorm.io/gorm"
)

// Storage представляет централизованное хранилище для всех репозиториев.
type Storage struct {
	Users     UsersRepository
	Cars      CarsRepository
	Orders    OrdersRepository
	Locations LocationsRepository
	Reviews   ReviewsRepository
}

// Интерфейсы для репозиториев

// UsersRepository определяет методы для работы с пользователями.
type UsersRepository interface {
	Create(email, password string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}

// CarsRepository определяет методы для работы с автомобилями.
type CarsRepository interface {
	Create(car *models.Car) (uint, error)
	GetByID(id uint) (*models.Car, error)
	GetBrands() ([]string, error)
	GetModelsByBrand(brand string) ([]string, error)
	Update(car *models.Car) error
	Delete(id uint) error
	FilterByStatus(status string) ([]models.Car, error)
	FilterByLocation(locationID uint) ([]models.Car, error)
	Filter(req models.SearchCarRequest) ([]models.Car, error)
	Search(query string) ([]models.Car, error)
	SortByPrice(ascending bool) ([]models.Car, error)
}

// OrdersRepository определяет методы для работы с заказами.
type OrdersRepository interface {
	Create(order *models.Order) (uint, error)
	GetByID(id uint) (*models.Order, error)
	GetAll() ([]models.Order, error)
	Update(order *models.Order) error
	Delete(id uint) error
	FilterByStatus(status string) ([]models.Order, error)
	FilterByUserID(userID uint) ([]models.Order, error)
	FilterByCarID(carID uint) ([]models.Order, error)
}

// LocationsRepository определяет методы для работы с локациями.
type LocationsRepository interface {
	Create(location *models.Location) (uint, error)
	GetByID(id uint) (*models.Location, error)
	GetAll() ([]models.Location, error)
	Update(location *models.Location) error
	Delete(id uint) error
	FilterByCity(city string) ([]models.Location, error)
	FilterByCountry(country string) ([]models.Location, error)
}

// ReviewsRepository определяет методы для работы с отзывами.
type ReviewsRepository interface {
	Create(review *models.Review) (uint, error)
	GetByID(id uint) (*models.Review, error)
	GetAll() ([]models.Review, error)
	Update(review *models.Review) error
	Delete(id uint) error
	FilterByUserID(userID uint) ([]models.Review, error)
	FilterByCarID(carID uint) ([]models.Review, error)
	FilterByRating(rating int) ([]models.Review, error)
}

// New создает новый экземпляр Storage с инициализированными репозиториями.
func New(db *gorm.DB) *Storage {
	return &Storage{
		Users:     NewUsersRepo(db),
		Cars:      NewCarsRepo(db),
		Orders:    NewOrdersRepo(db),
		Locations: NewLocationsRepo(db),
		Reviews:   NewReviewsRepo(db),
	}
}
