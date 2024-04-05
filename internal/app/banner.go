package banner

import (
	tr "banner/internal/transport"
	"log"
	"net/http"
)

func StartServer() {
	tr.ConfigTransport()
	log.Fatal(http.ListenAndServe(":8080", tr.Router))
}
