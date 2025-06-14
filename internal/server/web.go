package server

import (
	"ISIT/internal/config"
	"ISIT/internal/storage"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log/slog"
)

type Server struct {
	app     *echo.Echo
	URL     string
	logger  *slog.Logger
	Storage *storage.Storage
	JWT     config.JWT
}

func New(srvCfg config.Server, Jwt config.JWT, logger *slog.Logger, storage *storage.Storage) (*Server, error) {
	e := echo.New()
	server := Server{
		app:     e,
		URL:     srvCfg.URL,
		logger:  logger,
		Storage: storage,
		JWT:     Jwt,
	}
	e.HideBanner = true
	e.Logger.SetOutput(io.Discard)

	// Middleware для всего приложения
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	// Статические файлы
	e.Static("/static", "static")

	//// Обработчик для base.html (использовать для MPA версии)
	//e.GET("/", func(c echo.Context) error {
	//	// Проверяем существование файла
	//	filePath := filepath.Join("static/index", "base.html")
	//	return c.File(filePath)
	//})

	// Создаем middleware для проверки JWT
	m := NewMiddleware(Jwt, logger)

	// Регистрация обработчиков
	server.RegisterHandlers(m)

	return &server, nil
}

func (s *Server) Serve() error {
	s.logger.Info("HTTP server started", slog.String("url", s.URL))
	return fmt.Errorf("server error: %w", s.app.Start(s.URL))
}
