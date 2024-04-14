package banner

import (
	m "banner/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

func GetAllBanners(featureID, tagID, limit, offset int, db *sql.DB) ([]m.ListOfBanners, error) {
	query := `
        SELECT b.banner_id, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at, bt.tag_id
        FROM banners b
        LEFT JOIN banner_tags bt ON b.banner_id = bt.banner_id
        WHERE ($1 = 0 OR b.feature_id = $1)
        AND ($2 = 0 OR bt.tag_id = $2)
    `

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	query += fmt.Sprintf(" OFFSET %d", offset)

	rows, err := db.Query(query, featureID, tagID)
	if err != nil {
		return nil, fmt.Errorf("error querying banners: %v", err)
	}
	defer rows.Close()

	banners := make([]m.ListOfBanners, 0)
	for rows.Next() {
		var banner m.ListOfBanners
		var tagID sql.NullInt64
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&banner.BannerID, &banner.FeatureID, &banner.Title, &banner.Text, &banner.URL, &banner.IsActive,
			&createdAt, &updatedAt, &tagID); err != nil {
			return nil, fmt.Errorf("error scanning banner: %v", err)
		}

		banner.CreatedAt = createdAt
		banner.UpdatedAt = updatedAt

		if tagID.Valid {
			banner.TagIDs = append(banner.TagIDs, int(tagID.Int64))
		}

		banners = append(banners, banner)
	}

	return banners, nil
}

func FetchBannerFromDB(db *sql.DB, tagID, featureID int) (*m.ResponseBanner, error) {
	if db == nil {
		log.Println("FetchBannerFromDB called with nil database connection")
		return nil, errors.New("database connection is nil")
	}

	query := `
        SELECT b.title, b.text, b.url, b.is_active
        FROM banners b
        JOIN banner_tags t ON b.banner_id = t.banner_id
        WHERE t.tag_id = $1 AND b.feature_id = $2;
    `
	banner := &m.ResponseBanner{
		Content: &m.Content{},
	}

	if err := db.QueryRow(query, tagID, featureID).Scan(
		&banner.Content.Title,
		&banner.Content.Text,
		&banner.Content.URL,
		&banner.IsActive,
	); err != nil {
		log.Printf("Error fetching banner from DB: %v", err)
		return nil, err
	}
	log.Println("Successfully fetched banner from DB")
	return banner, nil
}
