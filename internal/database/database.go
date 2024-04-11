package banner

import (
	m "banner/internal/models"
	"database/sql"
	"fmt"
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

func FindUserByToken(token string) (*m.User, error) {
	user := m.User{}
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

func UpdateBannerInDB(bannerID int, requestData struct {
	TagIDs    []int    `json:"tag_ids,omitempty"`
	FeatureID *int     `json:"feature_id,omitempty"`
	Content   m.Banner `json:"content,omitempty"`
	IsActive  *bool    `json:"is_active,omitempty"`
}, db *sql.DB) error {
	// Формируем SQL запрос для обновления баннера.
	query := "UPDATE banners SET"
	args := make([]interface{}, 0)

	if requestData.FeatureID != nil {
		query += " feature_id = ?,"
		args = append(args, *requestData.FeatureID)
	}

	if len(requestData.TagIDs) > 0 {
		// Предполагается, что у вас есть отдельная таблица tag_banner, связывающая баннеры с тегами.
		// Вам нужно будет выполнить соответствующие действия для обновления связей с тегами.
		// Здесь приведен только пример формирования части SQL запроса.
		query += " WHERE banner_id = ?"
		args = append(args, bannerID)
	}

	if requestData.Content.Title != "" {
		query += " title = ?,"
		args = append(args, requestData.Content.Title)
	}
	if requestData.Content.Text != "" {
		query += " text = ?,"
		args = append(args, requestData.Content.Text)
	}
	if requestData.Content.Url != "" {
		query += " url = ?,"
		args = append(args, requestData.Content.Url)
	}

	if requestData.IsActive != nil {
		query += " is_active = ?,"
		args = append(args, *requestData.IsActive)
	}

	query += " WHERE banner_id = ?"
	args = append(args, bannerID)

	// Выполняем SQL запрос.
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баннера: %v", err)
	}

	return nil
}

func UpdateTagBannerInDB(bannerID int, tagIDs []int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM banner_tags WHERE banner_id = $1", bannerID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении существующих связей: %v", err)
	}

	// Добавляем новые связи.
	for _, tagID := range tagIDs {
		_, err := db.Exec("INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)", bannerID, tagID)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении связи: %v", err)
		}
	}

	return nil
}
