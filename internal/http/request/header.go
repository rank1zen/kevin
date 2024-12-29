package request

import "net/http"

func IsHTMX(r *http.Request) bool {
	value := r.Header.Get("HX-Request")
	return value == "true"
}
