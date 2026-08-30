package runtime

import "net/http"

func HandleReadyz(w http.ResponseWriter, r *http.Request) {
	// TODO: check deps are ready

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	_, _ = w.Write([]byte("ready"))
}
