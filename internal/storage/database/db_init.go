package banner

import (
	"database/sql"
	"log"
	"os"
)

var Db *sql.DB

func CreateDB() {
	file, err := os.ReadFile("migrations/storage.sql")
	if err != nil {
		panic(err)
	}

	log.Println("storage started to create")

	_, err = Db.Exec(string(file))
	if err != nil {
		log.Printf("Error executing storage update: %v", err)
		panic(err)
	}

	log.Println("storage created successfully")
}
