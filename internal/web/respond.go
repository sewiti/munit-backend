package web

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/apex/log"
	"github.com/go-sql-driver/mysql"
	"github.com/sewiti/munit-backend/internal/model"
)

var (
	errForbidden        = errors.New("403 Forbidden")
	errUnsupportedMedia = errors.New("415 Unsupported Media Type")
	errInternalError    = errors.New("500 Internal Server Error")
)

// respond responds with a JSON encoded body.
func respond(w http.ResponseWriter, body interface{}, code int) {
	if code == http.StatusNoContent || body == nil {
		w.WriteHeader(code)
		return
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		log.WithError(err).WithField("body", body).Error("unable to encode json")
	}
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
//	- errForbidden:          403
//	- model.ErrNotFound:     404
//	- sql.ErrNoRows:         404
//	- errUnsupportedContent: 415
//	- errInternalError:      500
func respondErr(w http.ResponseWriter, err error) {
	code := http.StatusBadRequest
	if sqlErr, ok := err.(*mysql.MySQLError); ok {
		log.WithError(sqlErr).Error("mysql error")
		respondInternalError(w)
		return
	}

	switch {
	case errors.Is(err, sql.ErrNoRows):
		err = model.ErrNotFound
		code = http.StatusNotFound

	case errors.Is(err, model.ErrNotFound):
		code = http.StatusNotFound

	case errors.Is(err, errForbidden):
		code = http.StatusForbidden

	case errors.Is(err, errUnsupportedMedia):
		code = http.StatusUnsupportedMediaType

	case errors.Is(err, errInternalError):
		respondInternalError(w)
		return
	}
	respondMsg(w, err.Error(), code)
}

// respondUnauthorized is a shorthand for
//	respondMsg(w, "401 Unauthorized", http.StatusUnauthorized)
func respondUnauthorized(w http.ResponseWriter) {
	respondMsg(w, "401 Unauthorized", http.StatusUnauthorized)
}

// respondInternalError is a shorthand for
//	respondErr(w, "500 Internal Server Error", http.StatusInternalServerError)
func respondInternalError(w http.ResponseWriter) {
	respondMsg(w, "500 Internal Server Error", http.StatusInternalServerError)
}
