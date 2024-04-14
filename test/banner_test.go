package banner

import (
	c "banner/internal/storage/database"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFindUserByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub storage connection", err)
	}
	defer db.Close()

	token := "some_token"
	rows := sqlmock.NewRows([]string{"user_id", "token", "is_admin"}).
		AddRow(1, token, true)

	mock.ExpectQuery("SELECT user_id, token, is_admin FROM users WHERE token = ?").WithArgs(token).WillReturnRows(rows)

	user, err := c.FindUserByToken(token, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if user == nil {
		t.Error("expected user, got nil")
		return
	}

	if user.UserID != 1 || user.Token != token || !user.IsAdmin {
		t.Errorf("unexpected user data: %+v", user)
	}
}
func TestGetAllBanners(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub storage connection", err)
	}
	defer db.Close()

	featureID, tagID, limit, offset := 20, 30, 10, 0
	rows := sqlmock.NewRows([]string{"banner_id", "feature_id", "title", "text", "url", "is_active", "created_at", "updated_at", "tag_id"}).
		AddRow(1, featureID, "Test Banner", "Test Text", "http://testurl.com", true, time.Now(), time.Now(), tagID)

	query := `SELECT b.banner_id, b.feature_id, b.title, b.text, b.url, b.is_active, b.created_at, b.updated_at, bt.tag_id
			  FROM banners b
			  LEFT JOIN banner_tags bt ON b.banner_id = bt.banner_id
			  WHERE ($1 = 0 OR b.feature_id = $1) AND ($2 = 0 OR bt.tag_id = $2) LIMIT $3 OFFSET COALESCE($4, 0)`
	mock.ExpectQuery(query).WithArgs(featureID, tagID, limit, offset).WillReturnRows(rows)

	banners, err := c.GetAllBanners(featureID, tagID, limit, offset, db)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if len(banners) != 1 {
		t.Errorf("expected one banner, got %d", len(banners))
		return
	}

	banner := banners[0]
	if banner.Title != "Test Banner" || banner.Text != "Test Text" || banner.Url != "http://testurl.com" || !banner.IsActive {
		t.Errorf("unexpected banner data: %+v", banner)
	}
}
