package banner

import (
	s "banner/internal/services"
	"github.com/gorilla/mux"
	"net/http"
)

func SetupRoutes(router *mux.Router, server *s.Instance) {
	router.Use(AuthMiddleware(server))
	router.HandleFunc("/user_banner", server.GetUserBanner).Methods("GET")

	router.HandleFunc("/banner", withAdminCheck(server, server.GetAllBanners)).Methods("GET")
	router.HandleFunc("/banner", withAdminCheck(server, server.CreateBanner)).Methods("POST")
	router.HandleFunc("/banner/{id}", withAdminCheck(server, server.UpdateBanner)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", withAdminCheck(server, server.DeleteBanner)).Methods("DELETE")
}

func withAdminCheck(server *s.Instance, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		adminMiddleware := AdminCheckMiddleware(server)
		adminMiddleware(handler).ServeHTTP(w, r)
	}
}
