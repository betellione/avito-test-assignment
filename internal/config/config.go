package banner

import (
	cache "banner/internal/storage/cache"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {
	envConfig()
}
func envConfig() {
	if err := godotenv.Load("configs/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	viper.AutomaticEnv()
}
func DBConfig() *sql.DB {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("POSTGRES_USER"),
		viper.GetString("POSTGRES_DB"),
		viper.GetString("POSTGRES_PASSWORD"),
		viper.GetString("DB_PORT"))

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return db
}

func RouterConfig() *mux.Router {
	return mux.NewRouter()
}

func RedisConfig() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST"),
		Password: "",
		DB:       0,
	})
	_, err := rdb.Ping(cache.Ctx).Result()
	if err != nil {
		log.Fatal(err)
	}
	if err := rdb.Ping(cache.Ctx).Err(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}
	return rdb
}
