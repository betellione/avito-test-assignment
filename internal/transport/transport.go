package banner

import (
	s "banner/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

var Router *mux.Router

func SetupRoutes(router *mux.Router) {
	router.Use(AuthMiddleware)
	router.HandleFunc("/user_banner", s.GetUserBanner).Methods("GET")

	router.HandleFunc("/banner", withAdminCheck(s.GetAllBanners)).Methods("GET")
	router.HandleFunc("/banner", withAdminCheck(s.CreateBanner)).Methods("POST")
	router.HandleFunc("/banner/{id}", withAdminCheck(s.UpdateBanner)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", withAdminCheck(s.DeleteBanner)).Methods("DELETE")
}

func withAdminCheck(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		AdminCheckMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		})).ServeHTTP(w, r)
	}
}
