package request

import (
	"net/http"
	"strconv"
)

// QueryIntParam returns a query string parameter as integer.
func QueryIntParam(r *http.Request, param string, defaultValue int) int {
	value := r.URL.Query().Get(param)
	if value == "" {
		return defaultValue
	}

	val, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return defaultValue
	}

	if val < 0 {
		return defaultValue
	}

	return int(val)
}
