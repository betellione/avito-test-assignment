package banner

import (
	m "banner/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

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
