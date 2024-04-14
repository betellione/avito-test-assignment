package main

import (
	app "banner/internal/app"
	config "banner/internal/config"
	context "banner/internal/storage/database"
)

func main() {
	// TODO хеширование токена пользователя
	// TODO миграции
	// TODO логирование действий
	// TODO привести все ошибки к одному виду
	// TODO redis password
	config.InitConfig()
	context.CreateDB()
	app.StartServer(":8080")
}
