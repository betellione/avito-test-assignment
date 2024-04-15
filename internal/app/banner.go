package banner

import (
	s "banner/internal/services"
	tr "banner/internal/transport"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func StartServer(address string, router *mux.Router, s *s.Instance) {
	tr.SetupRoutes(router, s)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
