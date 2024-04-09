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

func FindUserByToken(token string) (*model.User, error) {
	user := model.User{}
	row := Db.QueryRow("SELECT user_id, token, is_admin FROM users WHERE token = $1", token)
	err := row.Scan(&user.UserID, &user.Token, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteBannerFromDB(bannerID int) error {
	query := `DELETE FROM banners WHERE banner_id = $1`

	result, err := Db.Exec(query, bannerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
