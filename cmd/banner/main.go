package main

import (
	app "banner/internal/app"
	config "banner/internal/config"
	context "banner/internal/database"
)

func main() {
	// TODO хеширование токена пользователя
	config.InitConfig()
	context.CreateOrUpdateDB()
	app.StartServer()
}
