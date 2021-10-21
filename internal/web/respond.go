package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sewiti/munit-backend/internal/model"
)

// respond responds with a JSON encoded body.
func respond(w http.ResponseWriter, body interface{}, code int) {
	if code == http.StatusNoContent || body == nil {
		w.WriteHeader(code)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

// respondMsg responds with a JSON:
//	{"message": "..."}
func respondMsg(w http.ResponseWriter, msg string, code int) {
	body := struct {
		Msg string `json:"message"`
	}{msg}
	respond(w, body, code)
}

// respondOK is a shorthand for
//	respond(w, body, http.StatusOK)
func respondOK(w http.ResponseWriter, body interface{}) {
	respond(w, body, http.StatusOK)
}

// respondErr responds based on error type.
//	model.ErrNotFound:     http.StatusNotFound,
//	ErrUnsupportedContent: http.StatusUnsupportedMediaType,
func respondErr(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	switch {
	case errors.Is(err, model.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, ErrUnsupportedContent):
		code = http.StatusUnsupportedMediaType
	}
	respondMsg(w, err.Error(), code)
}

// respondInternalError is a shorthand for
//	respondErr(w, "500 Internal Server Error", http.StatusInternalServerError)
func respondInternalError(w http.ResponseWriter) {
	respondMsg(w, "500 Internal Server Error", http.StatusInternalServerError)
}
