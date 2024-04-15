package banner

import (
	m "banner/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"time"
)

func GetAllBanners(featureID, tagID, limit, offset int, db *sql.DB) ([]m.ListOfBanners, error) {
	query := `
        WITH FilteredBanners AS (
            SELECT DISTINCT b.banner_id
            FROM banners b
            JOIN banner_tags bt ON b.banner_id = bt.banner_id
            WHERE ($1 = 0 OR b.feature_id = $1)
            AND ($2 = 0 OR bt.tag_id = $2)
        )
        SELECT b.banner_id, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at,
               array_agg(bt.tag_id) AS tags
        FROM banners b
        LEFT JOIN banner_tags bt ON b.banner_id = bt.banner_id
        WHERE b.banner_id IN (SELECT banner_id FROM FilteredBanners)
        GROUP BY b.banner_id, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at
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

	banners := []m.ListOfBanners{}
	for rows.Next() {
		var banner m.ListOfBanners
		var tags []int64
		var createdAt, updatedAt time.Time

		if err := rows.Scan(&banner.BannerID, &banner.FeatureID, &banner.Title, &banner.Text, &banner.URL, &banner.IsActive,
			&createdAt, &updatedAt, pq.Array(&tags)); err != nil {
			return nil, fmt.Errorf("error scanning banner: %v", err)
		}

		banner.CreatedAt = createdAt
		banner.UpdatedAt = updatedAt
		banner.TagIDs = make([]int, len(tags))
		for i, tagID := range tags {
			banner.TagIDs[i] = int(tagID)
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

func СheckBannerExists(bannerID int, db *sql.DB) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM banners WHERE banner_id = $1)", bannerID).Scan(&exists)
	return exists, err
}
