package banner

import (
	m "banner/internal/models"
	"database/sql"
)

func FindUserByToken(token string, db *sql.DB) (*m.User, error) {
	user := m.User{}
	row := db.QueryRow("SELECT user_id, token, is_admin FROM users WHERE token = $1", token)
	err := row.Scan(&user.UserID, &user.Token, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func IsAdminToken(token string, db *sql.DB) bool {
	user, err := FindUserByToken(token, db)
	if err != nil {
		return false
	}

	if !user.IsAdmin {
		return false
	}
	return true
}
