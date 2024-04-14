package banner

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func parseQueryInt(r *http.Request, key string) (int, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return 0, nil
	}
	return strconv.Atoi(valStr)
}

// parseJSONRequest parses the JSON body of a request.
func parseJSONRequest(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}
