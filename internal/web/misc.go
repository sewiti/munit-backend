package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sewiti/munit-backend/pkg/id"
)

func getIDs(r *http.Request, keys ...string) ([]id.ID, error) {
	vars := mux.Vars(r)
	if vars == nil {
		return nil, errors.New("vars is nil")
	}

	var err error
	ids := make([]id.ID, len(keys))
	for i, k := range keys {
		ids[i], err = id.Parse(vars[k])
		if err != nil {
			return nil, fmt.Errorf("%s: %w", k, err)
		}
	}
	return ids, nil
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
	return errUnsupportedMedia
}

func decodeJSONLimit(r *http.Request, v interface{}, limit int64) error {
	if err := assertJSON(r); err != nil {
		return err
	}
	return json.NewDecoder(io.LimitReader(r.Body, limit)).Decode(v)
}

func decodeJSON(r *http.Request, v interface{}) error {
	return decodeJSONLimit(r, v, defaultBodyLimit)
}
