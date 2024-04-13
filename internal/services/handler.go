package banner

import (
	context "banner/internal/database"
	m "banner/internal/models"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func GetUserBanner(w http.ResponseWriter, r *http.Request) {

}

func GetAllBanners(w http.ResponseWriter, r *http.Request) {
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

	banners, err := context.GetAllBanners(featureIDInt, tagIDInt, limitInt, offsetInt)
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

	bannerID, err := context.CreateBanner(requestData)
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

	err = context.UpdateBanner(bannerID, requestData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if err := context.DeleteBanner(bannerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Баннер для id не найден")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
