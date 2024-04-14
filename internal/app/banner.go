package banner

import (
	s "banner/internal/services"
	tr "banner/internal/transport"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartServer(address string, router *mux.Router, s *s.Instance) {
	tr.SetupRoutes(router, s)
	log.Fatal(http.ListenAndServe(address, router))
}
