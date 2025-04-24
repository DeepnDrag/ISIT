package models

// User представляет пользователя системы (клиента или администратора).
// User представляет пользователя системы (клиента или администратора).
type User struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Email        string `gorm:"unique;not null" json:"email"`
	PasswordHash string `gorm:"not null"`
	FirstName    string `gorm:"default:'не заполнено'" json:"first_name"`
	LastName     string `gorm:"default:'не заполнено'" json:"last_name"`
	PhoneNumber  string `gorm:"default:'не заполнено'" json:"phone_number"`
	Role         string `gorm:"not null;default:'customer'" json:"role"`
	CreatedAt    string `gorm:"not null" json:"created_at"`
	UpdatedAt    string `gorm:"not null" json:"updated_at"`

	// Связанные заказы и отзывы
	Orders  []Order  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Reviews []Review `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Car представляет автомобиль.
type Car struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Brand       string  `gorm:"not null" json:"brand"`
	Model       string  `gorm:"not null" json:"model"`
	Year        int     `gorm:"not null" json:"year"`
	Color       string  `gorm:"not null" json:"color"`
	Mileage     int     `gorm:"not null" json:"mileage"`
	PricePerDay float64 `gorm:"not null" json:"price_per_day"`
	Status      string  `gorm:"not null;default:available" json:"status"`
	LocationID  uint    `gorm:"not null" json:"location_id"`
	ImageURL    string  `gorm:"column:image_url" json:"image_url"`
	CreatedAt   string  `gorm:"not null" json:"created_at"`
	UpdatedAt   string  `gorm:"not null" json:"updated_at"`

	// Связанные заказы и отзывы
	Orders  []Order  `gorm:"foreignKey:CarID;constraint:OnDelete:CASCADE"`
	Reviews []Review `gorm:"foreignKey:CarID;constraint:OnDelete:CASCADE"`
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
	ID        uint    `gorm:"primaryKey; autoIncrement"`
	UserID    uint    `gorm:"not null"`
	CarID     uint    `gorm:"not null" json:"car_id"`
	StartDate string  `gorm:"not null" json:"start_date"` // Дата начала аренды
	EndDate   string  `gorm:"not null" json:"end_date"`   // Дата окончания аренды
	TotalCost float64 `gorm:"not null" json:"total_cost"` // Общая стоимость аренды
	Status    string  `gorm:"not null;default:pending"`   // 'pending', 'confirmed', 'completed', 'cancelled'
	CreatedAt string  `gorm:"not null"`
	UpdatedAt string  `gorm:"not null"`
}

// Review представляет отзыв пользователя о машине.
type Review struct {
	ID        uint   `gorm:"primaryKey; autoIncrement"`
	UserID    uint   `gorm:"not null"`
	CarID     uint   `gorm:"not null" json:"car_id"`
	Rating    int    `gorm:"not null" json:"rating"` // Оценка (например, от 1 до 5)
	Comment   string `gorm:"not null" json:"comment"`
	CreatedAt string `gorm:"not null" json:"created_at"`
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

type OrderWithCarRequest struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"-"`
	CarID     uint    `json:"car_id"`
	CarBrand  string  `json:"car_brand"`  // Марка автомобиля
	CarModel  string  `json:"car_model"`  // Модель автомобиля
	StartDate string  `json:"start_date"` // Дата начала аренды
	EndDate   string  `json:"end_date"`   // Дата окончания аренды
	TotalCost float64 `json:"total_cost"` // Общая стоимость аренды
	Status    string  `json:"status"`     // Статус заказа
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
