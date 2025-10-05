.PHONY: lint format logs test start start-d stop delete migrate delete-data

lint:
	@echo "==> Линтер"
	golangci-lint run

format:
	@echo "==> Форматирование"
	golangci-lint run --fix

logs:
	@echo "==> Логи"
	docker logs cloud-notes-server

test:
	@echo "==> Тесты"
	PROJECT_ROOT=$(pwd) go test -count=1 ./...

start:
	@echo "==> Запуск"
	docker compose up --build

start-d:
	@echo "==> Запуск в фоне"
	docker compose up --build -d

stop:
	@echo "==> Остановка"
	docker compose stop

delete:
	@echo "==> Удаление контейнеров"
	docker compose down

migrate:
	@echo "==> Миграции"
	docker compose up migrator

delete-data:
	@echo "==> Удаление данных"
	docker compose down -v
