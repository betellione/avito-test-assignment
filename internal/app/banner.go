package banner

import (
	tr "banner/internal/transport"
	"log"
	"net/http"
)

func StartServer(address string) {
	tr.SetupRoutes(tr.Router)
	log.Fatal(http.ListenAndServe(address, tr.Router))
}
