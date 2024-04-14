package banner

import (
	m "banner/internal/models"
	cache "banner/internal/storage/cache"
	context "banner/internal/storage/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type Instance struct {
	DB    *sql.DB
	Redis *redis.Client
}

func NewInstance(db *sql.DB, redis *redis.Client) *Instance {
	return &Instance{DB: db, Redis: redis}
}

func (s *Instance) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	tagID, featureID, err := parseBannerParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"

	banner, err := fetchBanner(tagID, featureID, useLastRevision, token, s.DB, s.Redis)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, banner.Content)
}

func fetchBanner(tagID, featureID int, useLastRevision bool, token string, db *sql.DB, r *redis.Client) (*m.ResponseBanner, error) {
	if !useLastRevision && !context.IsAdminToken(token, db) {
		if banner, err := cache.FetchBannerFromCache(r, tagID, featureID); err == nil {
			return banner, nil
		}
	}

	banner, err := context.FetchBannerFromDB(db, tagID, featureID)
	if err != nil {
		return nil, err
	}

	if !context.IsAdminToken(token, db) {
		cache.CacheBanner(r, tagID, featureID, banner)
	}

	return banner, nil
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func (s *Instance) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	bannerID, requestData, err := parseUpdateParams(r)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := context.UpdateBanner(bannerID, requestData, s.DB); err != nil {
		http.Error(w, "Failed to update banner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Instance) CreateBanner(w http.ResponseWriter, r *http.Request) {
	var requestData m.RequestData
	if err := parseJSONRequest(r, &requestData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.FeatureID == nil || requestData.Content.Title == "" || requestData.Content.Text == "" ||
		requestData.Content.URL == "" || requestData.IsActive == nil || len(requestData.TagIDs) == 0 {
		http.Error(w, "Insufficient data to create a banner", http.StatusBadRequest)
		return
	}

	bannerID, err := context.CreateBanner(requestData, s.DB)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"bannerId": bannerID})
}

func (s *Instance) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}

	if err := context.DeleteBanner(bannerID, s.DB); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Banner not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Instance) GetAllBanners(w http.ResponseWriter, r *http.Request) {
	featureID, tagID, limit, offset, err := parseListParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banners, err := context.GetAllBanners(featureID, tagID, limit, offset, s.DB)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, banners)
}
