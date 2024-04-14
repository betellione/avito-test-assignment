package banner

import (
	"database/sql"
	"log"
	"os"
)

func CreateDB(db *sql.DB) {
	file, err := os.ReadFile("migrations/database.sql")
	if err != nil {
		panic(err)
	}

	log.Println("storage started to create")

	_, err = db.Exec(string(file))
	if err != nil {
		log.Printf("Error executing storage update: %v", err)
		panic(err)
	}

	log.Println("storage created successfully")
}
