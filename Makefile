APP_NAME := server
CONFIG := ./config/local.yaml
CMD_PATH := ./cmd/server
DB_URL := postgres://admin:1234@localhost:5433/winter_db?sslmode=disable
MIGRATIONS_DIR := migrations

run:
	@echo "Запуск приложения..."
	go run $(CMD_PATH)/main.go -config=$(CONFIG)

migrate-create:
	@echo "Создание миграции: $(name)"
	@if [ -z "$(name)" ]; then \
		echo "Ошибка: нужно указать имя (пример: make migrate-create name=create_users_table)"; \
		exit 1; \
	fi
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

migrate-up:
	@echo "Применение миграций..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" up

migrate-down:
	@echo "Откат последней миграции..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down 1

migrate-reset:
	@echo "Полный откат всех миграций..."
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" down

migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_URL)" version

clean:
	@echo "Очистка..."
	rm -rf bin

.PHONY: build run test migrate-create migrate-up migrate-down migrate-reset migrate-version clean
