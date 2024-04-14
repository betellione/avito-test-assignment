BINARY_NAME=myapp
DOCKER_IMAGE_NAME=myapp-image

GO=go
DOCKER=docker
CP=cp

all: test build

build: env_setup
	$(GO) build -o $(BINARY_NAME) .

test: env_setup
	$(GO) test -v ./...

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

run: env_setup
	$(GO) run .

docker-build: test
	$(DOCKER) build -t $(DOCKER_IMAGE_NAME) .

docker-run: docker-build
	$(DOCKER) run --rm -p 8080:8080 $(DOCKER_IMAGE_NAME)

env_setup:
	$(CP) configs/example.env configs/.env

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
