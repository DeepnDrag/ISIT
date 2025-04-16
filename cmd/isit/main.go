package main

import (
	"ISIT/internal/config"
	"ISIT/internal/database"
	"ISIT/internal/logger"
	"ISIT/internal/server"
	"ISIT/internal/storage"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	config, err := config.New("config.yaml")
	if err != nil {
		return fmt.Errorf("new config error: %w", err)
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		return fmt.Errorf("new logger error: %w", err)
	}

	db, err := database.Connection(config.DB)
	if err != nil {
		return err
	}
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			logger.Error(err.Error())
			return
		}

		err = sqlDB.Close()
		if err != nil {
			logger.Error(err.Error())
			return
		}
	}()

	err = database.Migrations(db)
	if err != nil {
		logger.Error("database run migrations", err.Error())
		return err
	}

	st := storage.New(db)

	server, err := server.New(config.Server, config.JWT, logger, st)
	if err != nil {
		return err
	}

	return server.Serve()
}
