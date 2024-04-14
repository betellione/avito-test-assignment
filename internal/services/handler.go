package banner

import (
	m "banner/internal/models"
	cache "banner/internal/storage/cache"
	context "banner/internal/storage/database"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetUserBanner(w http.ResponseWriter, r *http.Request) {
	tagID, featureID, err := parseBannerParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := r.Header.Get("token")
	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"

	banner, err := fetchBanner(tagID, featureID, useLastRevision, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, banner)
}

func fetchBanner(tagID, featureID int, useLastRevision bool, token string) (*m.ResponseBanner, error) {
	if !useLastRevision && !context.IsAdminToken(token, context.Db) {
		if banner, err := cache.FetchBannerFromCache(cache.RedisClient, tagID, featureID); err == nil {
			return banner, nil
		}
	}
	return context.FetchBannerFromDB(context.Db, tagID, featureID)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func UpdateBanner(w http.ResponseWriter, r *http.Request) {
	bannerID, requestData, err := parseUpdateParams(r)
	if err != nil {
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	if err := context.UpdateBanner(bannerID, requestData, context.Db); err != nil {
		http.Error(w, "Failed to update banner", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func parseUpdateParams(r *http.Request) (int, m.RequestData, error) {
	idStr := r.URL.Path[len("/banner/"):]
	bannerID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, m.RequestData{}, err
	}

	var requestData m.RequestData
	if err = parseJSONRequest(r, &requestData); err != nil {
		return 0, m.RequestData{}, err
	}

	return bannerID, requestData, nil
}

func CreateBanner(w http.ResponseWriter, r *http.Request) {
	var requestData m.RequestData
	if err := parseJSONRequest(r, &requestData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	if requestData.FeatureID == nil || requestData.Content.Title == "" || requestData.Content.Text == "" ||
		requestData.Content.Url == "" || requestData.IsActive == nil || len(requestData.TagIDs) == 0 {
		http.Error(w, "Insufficient data to create a banner", http.StatusBadRequest)
		return
	}

	bannerID, err := context.CreateBanner(requestData, context.Db)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"bannerId": bannerID})
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerID, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "Invalid banner ID", http.StatusBadRequest)
		return
	}

	if err := context.DeleteBanner(bannerID, context.Db); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Banner not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetAllBanners(w http.ResponseWriter, r *http.Request) {
	featureID, tagID, limit, offset, err := parseListParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banners, err := context.GetAllBanners(featureID, tagID, limit, offset, context.Db)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, banners)
}
