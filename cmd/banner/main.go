package main

import (
	app "banner/internal/app"
	config "banner/internal/config"
	s "banner/internal/services"
	context "banner/internal/storage/database"
)

func main() {
	// TODO хеширование токена пользователя
	// TODO миграции
	// TODO логирование действий
	// TODO привести все ошибки к одному виду
	// TODO redis password
	config.InitConfig()
	db := config.DbConfig()
	context.CreateDB(db)
	app.StartServer(":8080", config.RouterConfig(), s.NewInstance(db, config.RedisConfig()))
}
