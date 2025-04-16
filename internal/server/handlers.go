package server

import (
	"ISIT/internal/models"
	"ISIT/internal/utils"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (s *Server) RegisterHandlers(m *Middleware) {
	app := s.app

	authGroup := app.Group("/auth")
	apiGroup := app.Group("/api")

	// Pages
	apiGroup.GET("/profile/page", s.handleProfilePage)
	apiGroup.GET("/search/page", s.handleSearchPage)
	apiGroup.GET("/car/page", s.handleCarPage)
	authGroup.GET("/login/page", s.handleAuthPage)

	// Users
	authGroup.POST("/login", s.Authorize)
	apiGroup.GET("/user/get", s.GetUserInfoByEmail, m.AccessLog()) // Получить пользователя по token
	apiGroup.POST("/user/update", s.UpdateUser, m.AccessLog())     // Обновить пользователя
	apiGroup.DELETE("/users/:id", s.DeleteUser)                    // Удалить пользователя

	// Cars
	apiGroup.GET("/cars/:id", s.GetCarDetails, m.AccessLog())
	apiGroup.POST("/car/add", s.CreateCar, m.AccessLog())         // Создать автомобиль
	apiGroup.GET("/cars/brands", s.ListCarsBrands, m.AccessLog()) // Список всех автомобилей
	apiGroup.GET("/cars/models", s.ListCarsModels, m.AccessLog())
	apiGroup.DELETE("/cars/:id", s.DeleteCar)  // Удалить автомобиль
	apiGroup.GET("/cars/filter", s.FilterCars) // Фильтрация автомобилей

	// Locations
	apiGroup.POST("/locations", s.CreateLocation)       // Создать локацию
	apiGroup.GET("/locations/:id", s.GetLocation)       // Получить локацию по ID
	apiGroup.GET("/locations", s.ListLocations)         // Список всех локаций
	apiGroup.PUT("/locations/:id", s.UpdateLocation)    // Обновить локацию
	apiGroup.DELETE("/locations/:id", s.DeleteLocation) // Удалить локацию

	// Orders
	apiGroup.POST("/orders", s.CreateOrder)       // Создать заказ
	apiGroup.GET("/orders/:id", s.GetOrder)       // Получить заказ по ID
	apiGroup.GET("/orders", s.ListOrders)         // Список всех заказов
	apiGroup.PUT("/orders/:id", s.UpdateOrder)    // Обновить заказ
	apiGroup.DELETE("/orders/:id", s.DeleteOrder) // Удалить заказ

	// Reviews
	apiGroup.POST("/reviews", s.CreateReview)       // Создать отзыв
	apiGroup.GET("/reviews/:id", s.GetReview)       // Получить отзыв по ID
	apiGroup.GET("/reviews", s.ListReviews)         // Список всех отзывов
	apiGroup.PUT("/reviews/:id", s.UpdateReview)    // Обновить отзыв
	apiGroup.DELETE("/reviews/:id", s.DeleteReview) // Удалить отзыв

	apiGroup.Use(m.AccessLog()) // Применяем middleware для проверки JWT

}

func (s *Server) handleCarPage(c echo.Context) error {
	return c.File("static/index/car.html")
}

// Обработчик для страницы ЛК
func (s *Server) handleProfilePage(c echo.Context) error {
	// Проверка JWT уже выполнена middleware, так что пользователь авторизован

	// Пример: Возвращаем статический файл profile.html
	return c.File("static/index/profile.html")
}

func (s *Server) handleSearchPage(c echo.Context) error {
	// Проверка JWT уже выполнена middleware, так что пользователь авторизован

	// Пример: Возвращаем статический файл profile.html
	return c.File("static/index/search.html")
}

// Обработчик страницы входа
func (s *Server) handleAuthPage(c echo.Context) error {
	return c.File("static/index/login.html")
}

func (s *Server) Authorize(c echo.Context) error {
	var req models.AuthorizeUserRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, err := s.Storage.Users.GetByEmail(req.Email)

	if err != nil {
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Password hashing error"})
		}
		createdUser, createErr := s.Storage.Users.Create(req.Email, hashedPassword)
		if createErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]error{"error": createErr})
		}
		user = createdUser
	} else {
		if !utils.CheckPassword(req.Password, user.PasswordHash) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
	}

	token, err := utils.GenerateToken(user.Email, s.JWT.SecretKey)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func (s *Server) GetCarDetails(c echo.Context) error {
	carID := c.Param("id")
	id, err := strconv.Atoi(carID)
	if err != nil {
		log.Println("aaaaaaaaaaa", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid car ID"})
	}

	// Получаем автомобиль по ID
	car, err := s.Storage.Cars.GetByID(uint(id))
	if err != nil {
		log.Println("bbbbbbb", err)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "car not found"})
	}
	log.Println("carcarcarcarcaracr", car)
	return c.JSON(http.StatusOK, car)
}

// FilterCars возвращает отфильтрованный список автомобилей
func (s *Server) FilterCars(c echo.Context) error {
	var req models.SearchCarRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Вызываем метод из storage для фильтрации
	// Вызываем сервисный метод
	cars, err := s.Storage.Cars.Filter(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to search cars"})
	}

	return c.JSON(http.StatusOK, cars)
}

func (s *Server) GetUserInfoByEmail(c echo.Context) error {
	// Извлекаем email из контекста, который был добавлен middleware при валидации токена
	email, ok := c.Get("email").(string)
	log.Println("Email from context:", email)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing or invalid token"})
	}

	// Получаем пользователя по email
	user, err := s.Storage.Users.GetByEmail(email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "User not found"})
	}
	log.Println(user.CreatedAt, user.UpdatedAt)
	// Возвращаем данные пользователя
	return c.JSON(http.StatusOK, map[string]interface{}{
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
		"created_at":   user.CreatedAt,
		"updated_at":   user.UpdatedAt,
		"first_name":   user.FirstName,
		"last_name":    user.LastName,
	})
}

// ListUsers возвращает список всех пользователей
func (s *Server) ListUsers(c echo.Context) error {
	users, err := s.Storage.Users.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch users"})
	}
	return c.JSON(http.StatusOK, users)
}

func (s *Server) UpdateUser(c echo.Context) error {
	// Получаем email из контекста (добавлен middleware)
	email, ok := c.Get("email").(string)
	slog.Info(email)
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing email in token"})
	}
	// Получаем пользователя по email
	user, err := s.Storage.Users.GetByEmail(email)
	log.Println(err)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}

	// Создаем структуру для обновления
	updatedUser := &models.User{}

	// Связываем данные из запроса с объектом updatedUser
	if err := c.Bind(updatedUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	// Обновляем только изменённые поля
	if updatedUser.FirstName != "" {
		user.FirstName = updatedUser.FirstName
	}
	if updatedUser.LastName != "" {
		user.LastName = updatedUser.LastName
	}
	if updatedUser.PhoneNumber != "" {
		user.PhoneNumber = updatedUser.PhoneNumber
	}
	if updatedUser.Role != "" {
		user.Role = updatedUser.Role
	}

	// Обновляем дату последнего изменения
	user.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	// Обновляем пользователя в базе
	if err := s.Storage.Users.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update user"})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser удаляет пользователя
func (s *Server) DeleteUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := s.Storage.Users.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete user"})
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (s *Server) CreateCar(c echo.Context) error {
	// Получаем текстовые данные из формы
	brand := c.FormValue("brand")
	model := c.FormValue("model")
	yearStr := c.FormValue("year")
	color := c.FormValue("color")
	mileageStr := c.FormValue("mileage")
	pricePerDayStr := c.FormValue("price_per_day")
	status := c.FormValue("status")
	locationIDStr := c.FormValue("location_id")

	// Преобразуем строки в нужные типы
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid year"})
	}

	mileage, err := strconv.Atoi(mileageStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid mileage"})
	}

	pricePerDay, err := strconv.ParseFloat(pricePerDayStr, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid price_per_day"})
	}

	locationID, err := strconv.Atoi(locationIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid location_id"})
	}

	// Получаем файл из формы
	file, header, err := c.Request().FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image is required"})
	}
	defer file.Close()

	// Генерируем уникальное имя файла
	extension := filepath.Ext(header.Filename)
	fileName := fmt.Sprintf("%s%s", time.Now().Format("20060102150405"), extension)
	filePath := filepath.Join("static", "images", fileName)

	// Создаем директорию, если её нет
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create directory"})
	}

	// Сохраняем файл на диск
	dst, err := os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save image"})
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to copy image"})
	}

	// Создаем объект Car
	car := models.Car{
		Brand:       brand,
		Model:       model,
		Year:        year,
		Color:       color,
		Mileage:     mileage,
		PricePerDay: pricePerDay,
		Status:      status,
		LocationID:  uint(locationID),
		ImageURL:    "/" + filePath,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}

	// Логируем полученные данные
	slog.Info("Received car data", "car", car)

	// Создаем запись в базе данных
	id, err := s.Storage.Cars.Create(&car)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create car"})
	}

	// Возвращаем успешный ответ
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "car created successfully",
		"car":     car,
	})
}

// GetCar получает автомобиль по ID
func (s *Server) GetCar(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	car, err := s.Storage.Cars.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "car not found"})
	}
	return c.JSON(http.StatusOK, car)
}

// ListCarsBrands возвращает список всех марок
func (s *Server) ListCarsBrands(c echo.Context) error {
	cars, err := s.Storage.Cars.GetBrands()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch cars"})
	}
	return c.JSON(http.StatusOK, cars)
}

// ListCarsModels по марке автомобиля выводит список доступных моделей
func (s *Server) ListCarsModels(c echo.Context) error {
	// Получаем название марки из параметров запроса
	brand := c.QueryParam("brand")
	if brand == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "brand parameter is required"})
	}

	// Получаем список моделей для указанной марки
	model, err := s.Storage.Cars.GetModelsByBrand(brand)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch models"})
	}

	// Возвращаем список моделей
	return c.JSON(http.StatusOK, map[string]interface{}{
		"models": model,
	})
}

// DeleteCar удаляет автомобиль
func (s *Server) DeleteCar(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := s.Storage.Cars.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete car"})
	}
	return c.JSON(http.StatusNoContent, nil)
}

// CreateLocation создаёт новую локацию
func (s *Server) CreateLocation(c echo.Context) error {
	var location models.Location
	if err := c.Bind(&location); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	id, err := s.Storage.Locations.Create(&location)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create location"})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":       id,
		"message":  "location created successfully",
		"location": location,
	})
}

// GetLocation получает локацию по ID
func (s *Server) GetLocation(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	location, err := s.Storage.Locations.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "location not found"})
	}
	return c.JSON(http.StatusOK, location)
}

// ListLocations возвращает список всех локаций
func (s *Server) ListLocations(c echo.Context) error {
	locations, err := s.Storage.Locations.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch locations"})
	}
	return c.JSON(http.StatusOK, locations)
}

// UpdateLocation обновляет данные локации
func (s *Server) UpdateLocation(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	location, err := s.Storage.Locations.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "location not found"})
	}
	if err := c.Bind(location); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := s.Storage.Locations.Update(location); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update location"})
	}
	return c.JSON(http.StatusOK, location)
}

// DeleteLocation удаляет локацию
func (s *Server) DeleteLocation(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := s.Storage.Locations.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete location"})
	}
	return c.JSON(http.StatusNoContent, nil)
}

// CreateOrder создаёт новый заказ
func (s *Server) CreateOrder(c echo.Context) error {
	var order models.Order
	if err := c.Bind(&order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	id, err := s.Storage.Orders.Create(&order)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create order"})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "order created successfully",
		"order":   order,
	})
}

// GetOrder получает заказ по ID
func (s *Server) GetOrder(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	order, err := s.Storage.Orders.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order not found"})
	}
	return c.JSON(http.StatusOK, order)
}

// ListOrders возвращает список всех заказов
func (s *Server) ListOrders(c echo.Context) error {
	orders, err := s.Storage.Orders.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch orders"})
	}
	return c.JSON(http.StatusOK, orders)
}

// UpdateOrder обновляет данные заказа
func (s *Server) UpdateOrder(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	order, err := s.Storage.Orders.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order not found"})
	}
	if err := c.Bind(order); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := s.Storage.Orders.Update(order); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update order"})
	}
	return c.JSON(http.StatusOK, order)
}

// DeleteOrder удаляет заказ
func (s *Server) DeleteOrder(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := s.Storage.Orders.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete order"})
	}
	return c.JSON(http.StatusNoContent, nil)
}

// CreateReview создаёт новый отзыв
func (s *Server) CreateReview(c echo.Context) error {
	var review models.Review
	if err := c.Bind(&review); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	id, err := s.Storage.Reviews.Create(&review)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create review"})
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "review created successfully",
		"review":  review,
	})
}

// GetReview получает отзыв по ID
func (s *Server) GetReview(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	review, err := s.Storage.Reviews.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "review not found"})
	}
	return c.JSON(http.StatusOK, review)
}

// ListReviews возвращает список всех отзывов
func (s *Server) ListReviews(c echo.Context) error {
	reviews, err := s.Storage.Reviews.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to fetch reviews"})
	}
	return c.JSON(http.StatusOK, reviews)
}

// UpdateReview обновляет данные отзыва
func (s *Server) UpdateReview(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	review, err := s.Storage.Reviews.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "review not found"})
	}
	if err := c.Bind(review); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if err := s.Storage.Reviews.Update(review); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update review"})
	}
	return c.JSON(http.StatusOK, review)
}

// DeleteReview удаляет отзыв
func (s *Server) DeleteReview(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
	}
	if err := s.Storage.Reviews.Delete(uint(id)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to delete review"})
	}
	return c.JSON(http.StatusNoContent, nil)
}
