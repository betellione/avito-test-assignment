package banner

import "time"

type Banner struct {
	BannerID  int       `db:"banner_id"`
	FeatureID int       `db:"feature_id"`
	Title     string    `db:"title"`
	Text      string    `db:"text"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Tag struct {
	TagID int    `db:"tag_id"`
	Name  string `db:"name"`
}

type TagBanner struct {
	BannerID int `db:"banner_id"`
	TagID    int `db:"tag_id"`
}

type Feature struct {
	FeatureID int    `db:"feature_id"`
	Name      string `db:"name"`
}

type User struct {
	UserID  int    `db:"user_id"`
	Token   string `db:"token"`
	IsAdmin bool   `db:"is_admin"`
}
