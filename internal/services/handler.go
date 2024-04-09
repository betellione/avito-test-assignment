package banner

import (
	context "banner/internal/database"
	"database/sql"
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
