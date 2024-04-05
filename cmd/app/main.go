package main

import (
	config "banner/internal/config"
	db "banner/internal/database"
	tr "banner/internal/transport"
	"fmt"
)

func main() {
	config.InitConfig()
	db.CreateOrUpdateDB()
	tr.ConfigTransport()
	fmt.Println("hello")
}
