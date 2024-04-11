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

}

func CreateBanner(w http.ResponseWriter, r *http.Request) {

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

	// Парсим тело запроса в структуру BannerUpdateRequest.
	var requestData struct {
		TagIDs    []int    `json:"tag_ids,omitempty"`
		FeatureID *int     `json:"feature_id,omitempty"`
		Content   m.Banner `json:"content,omitempty"`
		IsActive  *bool    `json:"is_active,omitempty"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	// Обновляем данные баннера в базе данных.
	// Например, используя структуру Banner и SQL запрос UPDATE.

	// Отправляем ответ об успешном выполнении операции.
	w.WriteHeader(http.StatusOK)
}

func DeleteBanner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bannerID, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if err := context.DeleteBannerFromDB(bannerID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Баннер для id не найден")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
