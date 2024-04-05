package banner

import (
	"database/sql"
	"log"
	"os"
)

var Db *sql.DB

func CreateOrUpdateDB() {
	file, err := os.ReadFile("migrations/database.sql")
	if err != nil {
		panic(err)
	}

	log.Println("database started to update")

	_, err = Db.Exec(string(file))
	if err != nil {
		panic(err)
	}
	log.Println("database updated successfully")
}
