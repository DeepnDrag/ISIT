GO_BUILD_PATH ?= $(CURDIR)/bin
GO_BUILD_APP_PATH ?= $(GO_BUILD_PATH)/isit/

# Цели для кросс-компиляции
GOOS ?= linux
GOARCH ?= amd64
CGO ?= 0

# Цель для сборки
build:
	go env -w GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO)
	go build -o $(GO_BUILD_APP_PATH) ./cmd/isit

# Запуск docker compose
up:
	docker-compose up --build -d

# Остановка docker compose
down:
	@docker-compose down

# Очистка сгенерированных файлов
clean:
	@rm -rf $(GO_BUILD_PATH)