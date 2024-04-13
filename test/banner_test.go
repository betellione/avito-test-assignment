package banner

import (
	c "banner/internal/database"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestFindUserByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
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
