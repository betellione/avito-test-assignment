package main

import (
	app "banner/internal/app"
	config "banner/internal/config"
)

func main() {
	// TODO хеширование токена пользователя
	// TODO миграции
	// TODO логирование действий
	config.InitConfig()
	//context.CreateOrUpdateDB()
	app.StartServer()
}
