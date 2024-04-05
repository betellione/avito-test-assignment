package main

import (
	app "banner/internal/app"
	config "banner/internal/config"
	db "banner/internal/database"
)

func main() {
	config.InitConfig()
	db.CreateOrUpdateDB()
	app.StartServer()
}
