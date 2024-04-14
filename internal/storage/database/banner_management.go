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
		requestData.Content.URL, requestData.IsActive).Scan(&bannerID)
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
	var queryParts []string
	var args []interface{}

	if requestData.FeatureID != nil {
		queryParts = append(queryParts, fmt.Sprintf("feature_id = $%d", len(args)+1))
		args = append(args, *requestData.FeatureID)
	}
	if requestData.Content.Title != "" {
		queryParts = append(queryParts, fmt.Sprintf("title = $%d", len(args)+1))
		args = append(args, requestData.Content.Title)
	}
	if requestData.Content.Text != "" {
		queryParts = append(queryParts, fmt.Sprintf("text = $%d", len(args)+1))
		args = append(args, requestData.Content.Text)
	}
	if requestData.Content.URL != "" {
		queryParts = append(queryParts, fmt.Sprintf("url = $%d", len(args)+1))
		args = append(args, requestData.Content.URL)
	}
	if requestData.IsActive != nil {
		queryParts = append(queryParts, fmt.Sprintf("is_active = $%d", len(args)+1))
		args = append(args, *requestData.IsActive)
	}

	if len(queryParts) == 0 {
		return errors.New("no data to update")
	}

	fullQuery := fmt.Sprintf("UPDATE banners SET %s WHERE banner_id = $%d", strings.Join(queryParts, ", "), len(args)+1)
	args = append(args, bannerID)

	if len(requestData.TagIDs) > 0 {
		if err := UpdateTagBanner(bannerID, requestData.TagIDs, db); err != nil {
			return fmt.Errorf("failed to update banner tags: %v", err)
		}
	}

	if _, err := db.Exec(fullQuery, args...); err != nil {
		return fmt.Errorf("failed to update banner: %v", err)
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
