package banner

import (
	m "banner/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var Db *sql.DB

func CreateOrUpdateDB(db *sql.DB) {
	file, err := os.ReadFile("migrations/storage.sql")
	if err != nil {
		panic(err)
	}

	log.Println("storage started to update")

	_, err = db.Exec(string(file))
	if err != nil {
		panic(err)
	}
	log.Println("storage updated successfully")
}

func FindUserByToken(token string, db *sql.DB) (*m.User, error) {
	user := m.User{}
	row := db.QueryRow("SELECT user_id, token, is_admin FROM users WHERE token = $1", token)
	err := row.Scan(&user.UserID, &user.Token, &user.IsAdmin)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteBanner(bannerID int, db *sql.DB) error {
	query := `DELETE FROM banners WHERE banner_id = $1`

	result, err := db.Exec(query, bannerID)
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

func UpdateBanner(bannerID int, requestData m.RequestData, db *sql.DB) error {
	query := "UPDATE banners SET"
	args := make([]interface{}, 0)
	argCounter := 1

	if requestData.FeatureID != nil {
		query += " feature_id = $1,"
		args = append(args, *requestData.FeatureID)
		argCounter++
	}

	if len(requestData.TagIDs) > 0 {
		err := UpdateTagBanner(bannerID, requestData.TagIDs, db)
		if err != nil {
			return err
		}
	}

	if requestData.Content.Title != "" {
		query += " title = $" + fmt.Sprint(argCounter) + ","
		args = append(args, requestData.Content.Title)
		argCounter++
	}
	if requestData.Content.Text != "" {
		query += " text = $" + fmt.Sprint(argCounter) + ","
		args = append(args, requestData.Content.Text)
		argCounter++
	}
	if requestData.Content.Url != "" {
		query += " url = $" + fmt.Sprint(argCounter) + ","
		args = append(args, requestData.Content.Url)
		argCounter++
	}
	if requestData.IsActive != nil {
		query += " is_active = $" + fmt.Sprint(argCounter) + ","
		args = append(args, *requestData.IsActive)
		argCounter++
	}

	if len(args) == 0 {
		return errors.New("нет данных для обновления")
	}

	query = strings.TrimSuffix(query, ",")

	query += " WHERE banner_id = $" + fmt.Sprint(argCounter)
	args = append(args, bannerID)

	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении баннера: %v", err)
	}

	return nil
}

func UpdateTagBanner(bannerID int, tagIDs []int, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM banner_tags WHERE banner_id = $1", bannerID)
	if err != nil {
		return fmt.Errorf("ошибка при удалении существующих связей: %v", err)
	}

	for _, tagID := range tagIDs {
		_, err := db.Exec("INSERT INTO banner_tags (banner_id, tag_id) VALUES ($1, $2)", bannerID, tagID)
		if err != nil {
			return fmt.Errorf("ошибка при добавлении связи: %v", err)
		}
	}

	return nil
}

func CreateBanner(requestData m.RequestData, db *sql.DB) (int, error) {
	query := `
        INSERT INTO banners (feature_id, title, text, url, is_active)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING banner_id
    `
	var bannerID int
	err := db.QueryRow(query, requestData.FeatureID, requestData.Content.Title, requestData.Content.Text,
		requestData.Content.Url, requestData.IsActive).Scan(&bannerID)
	if err != nil {
		return 0, fmt.Errorf("ошибка при создании баннера: %v", err)
	}

	err = UpdateTagBanner(bannerID, requestData.TagIDs, db)
	if err != nil {
		_, deleteErr := db.Exec("DELETE FROM banners WHERE banner_id = $1", bannerID)
		if deleteErr != nil {
			return 0, fmt.Errorf("ошибка при удалении баннера после неудачного присвоения тегов: %v", deleteErr)
		}
		return 0, fmt.Errorf("ошибка при присвоении тегов баннеру: %v", err)
	}

	return bannerID, nil
}

func GetAllBanners(featureID, tagID, limit, offset int, db *sql.DB) ([]m.ListOfBanners, error) {
	query := `
        SELECT b.banner_id, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at, bt.tag_id
        FROM banners b
        LEFT JOIN banner_tags bt ON b.banner_id = bt.banner_id
        WHERE ($1 = 0 OR b.feature_id = $1)
        AND ($2 = 0 OR bt.tag_id = $2)
    `

	args := []interface{}{featureID, tagID}

	if limit > 0 {
		query += " LIMIT $3"
		args = append(args, limit)
	}
	query += " OFFSET COALESCE($4, 0)"
	args = append(args, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	banners := make(map[int]m.ListOfBanners)
	for rows.Next() {
		var banner m.ListOfBanners
		var tagID sql.NullInt64
		var createdAt, updatedAt time.Time

		err := rows.Scan(&banner.BannerID, &banner.FeatureID, &banner.Title, &banner.Text, &banner.Url, &banner.IsActive,
			&createdAt, &updatedAt, &tagID)
		if err != nil {
			return nil, err
		}

		banner.CreatedAt = createdAt
		banner.UpdatedAt = updatedAt

		if tagID.Valid {
			banner.TagIDs = append(banner.TagIDs, int(tagID.Int64))
		}

		banners[banner.BannerID] = banner
	}

	result := make([]m.ListOfBanners, 0, len(banners))
	for _, banner := range banners {
		result = append(result, banner)
	}

	return result, nil
}

func FetchBannerFromDB(db *sql.DB, tagID, featureID int) (*m.Banner, error) {
	query := `
        SELECT b.title, b.text, b.url
        FROM banners b
        JOIN banner_tags t ON b.banner_id = t.banner_id
        WHERE t.tag_id = $1 AND b.feature_id = $2;
    `
	row := db.QueryRow(query, tagID, featureID)
	banner := &m.Banner{}
	err := row.Scan(&banner.Title, &banner.Text, &banner.Url)
	if err != nil {
		return nil, err
	}
	return banner, nil
}
