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
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		http.Error(w, "Invalid tag_id", http.StatusBadRequest)
		return
	}

	featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
	if err != nil {
		http.Error(w, "Invalid feature_id", http.StatusBadRequest)
		return
	}

	useLastRevision := r.URL.Query().Get("use_last_revision") == "true"
	token := r.Header.Get("token")

	var banner *m.ResponseBanner
	var cacheHit bool

	if !useLastRevision && !context.IsAdminToken(token, context.Db) {
		banner, err = cache.FetchBannerFromCache(cache.RedisClient, tagID, featureID)
		if err == nil {
			cacheHit = true
		}
	}

	if !cacheHit {
		banner, err = context.FetchBannerFromDB(context.Db, tagID, featureID)
		if err != nil {
			http.Error(w, "Banner not found", http.StatusNotFound)
			return
		}
		if !useLastRevision && !context.IsAdminToken(token, context.Db) {
			cache.CacheBanner(cache.RedisClient, tagID, featureID, banner)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banner)
}

func GetAllBanners(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	featureID := r.URL.Query().Get("feature_id")
	tagID := r.URL.Query().Get("tag_id")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	var err error
	var featureIDInt, tagIDInt, limitInt, offsetInt int
	if featureID != "" {
		featureIDInt, err = strconv.Atoi(featureID)
		if err != nil {
			http.Error(w, "Некорректное значение для feature_id", http.StatusBadRequest)
			return
		}
	}
	if tagID != "" {
		tagIDInt, err = strconv.Atoi(tagID)
		if err != nil {
			http.Error(w, "Некорректное значение для tag_id", http.StatusBadRequest)
			return
		}
	}
	if limitStr != "" {
		limitInt, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Некорректное значение для limit", http.StatusBadRequest)
			return
		}
	}
	if offsetStr != "" {
		offsetInt, err = strconv.Atoi(offsetStr)
		if err != nil {
			http.Error(w, "Некорректное значение для offset", http.StatusBadRequest)
			return
		}
	}

	banners, err := context.GetAllBanners(featureIDInt, tagIDInt, limitInt, offsetInt, context.Db)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(banners)
}

func CreateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var requestData m.RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	if requestData.FeatureID == nil || requestData.Content.Title == "" || requestData.Content.Text == "" ||
		requestData.Content.Url == "" || requestData.IsActive == nil || len(requestData.TagIDs) == 0 {
		http.Error(w, "Недостаточно данных для создания баннера", http.StatusBadRequest)
		return
	}

	bannerID, err := context.CreateBanner(requestData, context.Db)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]int{"bannerId": bannerID})
}

func UpdateBanner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/banner/"):]
	bannerID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некорректный идентификатор баннера", http.StatusBadRequest)
		return
	}
	var requestData m.RequestData
	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	err = context.UpdateBanner(bannerID, requestData, context.Db)
	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	if err := context.DeleteBanner(bannerID, context.Db); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Баннер для id не найден", http.StatusNotFound)
		} else {
			http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
