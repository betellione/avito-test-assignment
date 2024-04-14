# Определение переменных
BINARY_NAME=myapp
DOCKER_IMAGE_NAME=myapp-image

# Путь к исполняемому файлу
GO=go
DOCKER=docker
CP=cp

# Команда по умолчанию при вызове make без аргументов
all: test build

# Сборка проекта
build: env_setup
	$(GO) build -o $(BINARY_NAME) .

# Запуск тестов
test: env_setup
	$(GO) test -v ./...

# Очистка
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

# Запуск
run: env_setup
	$(GO) run .

# Сборка Docker образа
docker-build: test
	$(DOCKER) build -t $(DOCKER_IMAGE_NAME) .

# Запуск Docker контейнера
docker-run: docker-build
	$(DOCKER) run --rm -p 8080:8080 $(DOCKER_IMAGE_NAME)

# Настройка среды
env_setup:
	$(CP) configs/example.env configs/.env

# Помощь
help:
	@echo "Makefile для проекта ${BINARY_NAME}"
	@echo
	@echo "Выберите одну из следующих команд:"
	@echo "  make build         Собрать проект"
	@echo "  make test          Запуск тестов"
	@echo "  make clean         Очистка"
	@echo "  make run           Запуск"
	@echo "  make docker-build  Сборка Docker образа"
	@echo "  make docker-run    Запуск Docker контейнера"
	@echo "  make help          Помощь"
