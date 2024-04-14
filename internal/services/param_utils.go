package banner

import (
	"encoding/json"
	"fmt"
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

func parseJSONRequest(r *http.Request, target interface{}) error {
	return json.NewDecoder(r.Body).Decode(target)
}

func parseBannerParams(r *http.Request) (int, int, error) {
	tagID, err := parseQueryInt(r, "tag_id")
	if err != nil {
		return 0, 0, err
	}

	featureID, err := parseQueryInt(r, "feature_id")
	if err != nil {
		return 0, 0, err
	}

	return tagID, featureID, nil
}
func parseListParams(r *http.Request) (int, int, int, int, error) {
	featureID, err := parseQueryInt(r, "feature_id")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid feature_id")
	}
	tagID, err := parseQueryInt(r, "tag_id")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid tag_id")
	}
	limit, err := parseQueryInt(r, "limit")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid limit value")
	}
	offset, err := parseQueryInt(r, "offset")
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("invalid offset value")
	}
	return featureID, tagID, limit, offset, nil
}
