package web

import (
	"encoding/json"
	"net/http"
)

func respond(w http.ResponseWriter, body interface{}, code int) {
	if code == http.StatusNoContent || body == nil {
		w.WriteHeader(code)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func respondErr(w http.ResponseWriter, msg string, code int) {
	body := struct {
		Message string `json:"message"`
	}{
		Message: msg,
	}
	respond(w, body, code)
}

// respondOK is a shorthand for
//	respond(w, body, http.StatusOK)
func respondOK(w http.ResponseWriter, body interface{}) {
	respond(w, body, http.StatusOK)
}

// respondNoContent is a shorthand for
//	respond(w, nil, http.StatusNoContent)
func respondNoContent(w http.ResponseWriter) {
	respond(w, nil, http.StatusNoContent)
}

// respondBadRequest is a shorthand for
//	respondErr(w, err.Error(), http.StatusBadRequest)
func respondBadRequest(w http.ResponseWriter, err error) {
	respondErr(w, err.Error(), http.StatusBadRequest)
}

// respondNotFound is a shorthand for
//	respondErr(w, "404 Not Found", http.StatusNotFound)
func respondNotFound(w http.ResponseWriter, err error) {
	respondErr(w, err.Error(), http.StatusNotFound)
}

// respondInternalError is a shorthand for
//	respondErr(w, "500 Internal Server Error", http.StatusInternalServerError)
func respondInternalError(w http.ResponseWriter) {
	respondErr(w, "500 Internal Server Error", http.StatusInternalServerError)
}
