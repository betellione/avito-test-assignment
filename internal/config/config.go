package banner

import (
	db "banner/internal/database"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {
	envConfig()
	dbConfig()
}
func envConfig() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	viper.AutomaticEnv()
}
func dbConfig() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_DB"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("DB_PORT"))

	db.Db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
}
