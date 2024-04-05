package main

import "fmt"
import config "banner/internal/config"
import db "banner/internal/database"

func main() {
	config.InitConfig()
	db.CreateOrUpdateDB()
	fmt.Println("hello")
}
