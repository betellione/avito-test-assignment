package banner

import (
	s "banner/internal/services"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var Router *mux.Router

func ConfigTransport() {
	Router.HandleFunc("/user_banner", s.GetUserBanner).Methods("GET")
	Router.HandleFunc("/banner", s.GetAllBanners).Methods("GET")
	Router.HandleFunc("/banner", s.CreateBanner).Methods("POST")
	Router.HandleFunc("/banner/{id}", s.UpdateBanner).Methods("PATCH")
	Router.HandleFunc("/banner/{id}", s.DeleteBanner).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", Router))
}