package server

import (
	"ISIT/internal/config"
	"ISIT/internal/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Middleware struct {
	logger *slog.Logger
	JWT    config.JWT
}

func NewMiddleware(cfg config.JWT, logger *slog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
		JWT:    cfg,
	}
}

func (m *Middleware) Register(router *echo.Echo) {
	router.Use(m.AccessLog())
}

func (m *Middleware) AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()

			requestID := uuid.New().String()
			c.Set("requestID", requestID)

			m.logger.Info("Request started",
				slog.String("RequestID", requestID),
				slog.String("IP", c.RealIP()),
				slog.String("URL", c.Request().URL.Path),
				slog.String("Method", c.Request().Method),
			)

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				m.logger.Warn("Authorization header missing",
					slog.String("RequestID", requestID))
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Отсутствует токен авторизации"})
			}

			log.Println(authHeader)

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "некорректный формат токена"})
			}

			log.Println(parts[1])

			tokenString := parts[1]

			log.Println("lololool", m.JWT.SecretKey)
			claims, err := utils.ValidateJWT(tokenString, m.JWT.SecretKey)
			log.Println("CLAIMS", claims, err)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "токен недействителен"})
			}

			log.Println(claims)

			// Устанавливаем email в контекст с ключом "email"
			c.Set("email", claims.Email)
			slog.Info("DEBUG: Email установлен в контекст", claims.Email)
			c.Set("userID", claims.UserID)
			slog.Info("DEBUG: UserID установлен в контекст", claims.UserID)

			// Передаем управление следующему обработчику
			handlerErr := next(c)

			// Вычисляем время отклика
			responseTime := time.Since(startTime)
			if handlerErr != nil {
				m.logger.Error("Request Failed",
					slog.String("RequestID", requestID),
					slog.String("Time spent", strconv.FormatInt(int64(responseTime), 10)),
					slog.String("Error", handlerErr.Error()))
			} else {
				m.logger.Info("Request done",
					slog.String("RequestID", requestID),
					slog.String("Time spent", strconv.FormatInt(int64(responseTime), 10)),
				)
			}
			return handlerErr
		}
	}
}
