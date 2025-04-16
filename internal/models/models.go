package models

// User представляет пользователя системы (клиента или администратора).
type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`         // Первичный ключ, автоматически увеличивается
	Email        string `gorm:"unique;not null" json:"email"`               // Email (уникальный, обязательный)
	PasswordHash string `gorm:"not null" json:"-"`                          // Хэш пароля (обязательный, исключен из JSON)
	FirstName    string `gorm:"default:'не заполнено'" json:"first_name"`   // Имя (может быть пустым)
	LastName     string `gorm:"default:'не заполнено'" json:"last_name"`    // Фамилия (может быть пустой)
	PhoneNumber  string `gorm:"default:'не заполнено'" json:"phone_number"` // Номер телефона (может быть пустым)
	Role         string `gorm:"not null;default:'customer'" json:"role"`    // Роль ('customer' или 'admin', по умолчанию 'customer')
	CreatedAt    string `gorm:"not null" json:"created_at"`                 // Время создания записи
	UpdatedAt    string `gorm:"not null" json:"updated_at"`                 // Время обновления записи
}

type Car struct {
	ID          uint    `gorm:"primaryKey" json:"id"`                     // Идентификатор
	Brand       string  `gorm:"not null" json:"brand"`                    // Производитель (например, Toyota)
	Model       string  `gorm:"not null" json:"model"`                    // Модель (например, Camry)
	Year        int     `gorm:"not null" json:"year"`                     // Год выпуска
	Color       string  `gorm:"not null" json:"color"`                    // Цвет
	Mileage     int     `gorm:"not null" json:"mileage"`                  // Пробег
	PricePerDay float64 `gorm:"not null" json:"price_per_day"`            // Стоимость аренды за день
	Status      string  `gorm:"not null;default:available" json:"status"` // Статус: 'available', 'rented', 'maintenance'
	LocationID  uint    `gorm:"not null" json:"location_id"`              // ID локации
	ImageURL    string  `gorm:"column:image_url" json:"image_url"`
	CreatedAt   string  `gorm:"not null" json:"created_at"` // Дата создания
	UpdatedAt   string  `gorm:"not null" json:"updated_at"` // Дата последнего обновления
}

// Location представляет локацию, где находятся автомобили.
type Location struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null"` // Название локации (например, "Москва, офис 1")
	Address   string `gorm:"not null"`
	City      string `gorm:"not null"`
	Country   string `gorm:"not null"`
	CreatedAt string `gorm:"not null"`
	UpdatedAt string `gorm:"not null"`
}

// Order представляет заказ на аренду автомобиля.
type Order struct {
	ID        uint    `gorm:"primaryKey"`
	UserID    uint    `gorm:"not null"`
	CarID     uint    `gorm:"not null"`
	StartDate string  `gorm:"not null"`                 // Дата начала аренды
	EndDate   string  `gorm:"not null"`                 // Дата окончания аренды
	TotalCost float64 `gorm:"not null"`                 // Общая стоимость аренды
	Status    string  `gorm:"not null;default:pending"` // 'pending', 'confirmed', 'completed', 'cancelled'
	CreatedAt string  `gorm:"not null"`
	UpdatedAt string  `gorm:"not null"`
}

// Review представляет отзыв пользователя о машине.
type Review struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	User      User
	CarID     uint `gorm:"not null"`
	Car       Car
	Rating    int    `gorm:"not null"` // Оценка (например, от 1 до 5)
	Comment   string `gorm:"not null"`
	CreatedAt string `gorm:"not null"`
	UpdatedAt string `gorm:"not null"`
}

type AuthorizeUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SearchCarRequest struct {
	Brand     string  `json:"brand" query:"brand"`           // Марка автомобиля
	Model     string  `json:"model" query:"model"`           // Модель автомобиля
	YearFrom  int     `json:"year_from" query:"year_from"`   // Год выпуска (от)
	YearTo    int     `json:"year_to" query:"year_to"`       // Год выпуска (до)
	MinPrice  float64 `json:"min_price" query:"min_price"`   // Минимальная цена за день
	MaxPrice  float64 `json:"max_price" query:"max_price"`   // Максимальная цена за день
	StartDate string  `json:"start_date" query:"start_date"` // Дата начала аренды (в формате YYYY-MM-DD)
	EndDate   string  `json:"end_date" query:"end_date"`     // Дата окончания аренды (в формате YYYY-MM-DD)
}
