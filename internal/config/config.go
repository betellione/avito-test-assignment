package banner

import (
	context "banner/internal/storage"
	cache "banner/internal/storage/cache"
	transport "banner/internal/transport"
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
	dbConfig()
	routerConfig()
	redisConfig()
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

	context.Db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func routerConfig() {
	transport.Router = mux.NewRouter()
}

func redisConfig() {
	cache.RedisClient = redis.NewClient(&redis.Options{
		Addr: viper.GetString("REDIS_HOST"),
		//Password: viper.GetString("REDIS_PASSWORD"),
		Password: "",
		DB:       0,
	})
	_, err := cache.RedisClient.Ping(cache.Ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

}
