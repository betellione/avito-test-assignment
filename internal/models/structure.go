package banner

import "time"

type RequestData struct {
	TagIDs    []int  `json:"tag_ids,omitempty"`
	FeatureID *int   `json:"feature_id,omitempty"`
	Content   Banner `json:"content,omitempty"`
	IsActive  *bool  `json:"is_active,omitempty"`
	Offset    *int   `json:"offset,omitempty"`
	Limit     *int   `json:"limit,omitempty"`
}

type ListOfBanners struct {
	BannerID  int
	FeatureID int
	Title     string
	Text      string
	Url       string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	TagIDs    []int
}

type Content struct {
	Title string `db:"title"`
	Text  string `db:"text"`
	Url   string `db:"url"`
}

type ResponseBanner struct {
	Content  *Banner `json:"content"`
	IsActive bool    `json:"is_active"`
}
