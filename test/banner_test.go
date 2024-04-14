package banner

import (
	s "banner/internal/services"
	c "banner/internal/storage/database"
	tr "banner/internal/transport"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
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

func mockServer(t *testing.T) (*s.Instance, sqlmock.Sqlmock, *miniredis.Miniredis) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when starting miniredis", err)
	}
	redisClient := NewRedisClient()

	return s.NewInstance(db, redisClient), mock, mr
}
func NewRedisClient() *redis.Client {
	redisAddr := viper.GetString("REDIS_HOST")

	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
	}

	return rdb
}

func TestGetUserBanner(t *testing.T) {
	instance, mock, _ := mockServer(t)

	defer instance.Db.Close()
	router := mux.NewRouter()

	tr.SetupRoutes(router, instance)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT user_id, token, is_admin FROM users WHERE token = $1")).
		WithArgs("test_token").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "token", "is_admin"}).AddRow(1, "test_token", true))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT b.title, b.text, b.url, b.is_active FROM banners b JOIN banner_tags t ON b.banner_id = t.banner_id WHERE t.tag_id = $1 AND b.feature_id = $2")).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"title", "text", "url", "is_active"}).AddRow("Sample Title", "Sample Text", "http://sample.url", true))

	req, _ := http.NewRequest("GET", "/user_banner?tag_id=1&feature_id=1", nil)
	req.Header.Set("token", "test_token")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "handler returned wrong status code")

	expected := `{"Title":"Sample Title","Text":"Sample Text","Url":"http://sample.url"}`
	assert.JSONEq(t, expected, rr.Body.String(), "handler returned unexpected body")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
