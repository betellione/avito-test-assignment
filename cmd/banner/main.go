package main

import (
	app "banner/internal/app"
	config "banner/internal/config"
)

func main() {
	// TODO хеширование токена пользователя
	// TODO миграции
	// TODO логирование действий
	// TODO привести все ошибки к одному виду
	// TODO redis password
	config.InitConfig()
	//context.CreateOrUpdateDB()
	app.StartServer(":8080")
}
