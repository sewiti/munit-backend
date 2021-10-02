package web

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getUUIDs(r *http.Request, keys ...string) ([]uuid.UUID, error) {
	vars := mux.Vars(r)
	var uuids []uuid.UUID
	for _, k := range keys {
		id, err := uuid.Parse(vars[k])
		if err != nil {
			return nil, err
		}
		uuids = append(uuids, id)
	}
	return uuids, nil
}

func assertJSON(r *http.Request) error {
	const contentJSON = "application/json"
	content := r.Header.Get("Content-Type")

	// application/json
	if content == contentJSON {
		return nil
	}
	// application/json; charset=utf-8
	if strings.HasPrefix(content, contentJSON+";") {
		return nil
	}
	return errors.New("unsupported content type, expected " + contentJSON)
}

func decodeJSONLimit(r *http.Request, v interface{}, limit int64) error {
	if err := assertJSON(r); err != nil {
		return err
	}
	return json.NewDecoder(io.LimitReader(r.Body, limit)).Decode(v)
}

func decodeJSON(r *http.Request, v interface{}) error {
	return decodeJSONLimit(r, v, bodyLimit)
}
