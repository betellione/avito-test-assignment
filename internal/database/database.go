package banner

import (
	model "banner/internal/models"
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

func FindUserByToken(db *sql.DB, token string) (*model.User, error) {
	user := model.User{}
	row := db.QueryRow("SELECT user_id, token, is_admin FROM users WHERE token = ?", token)
	err := row.Scan(&user.UserID, &user.Token, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
