package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sewiti/munit-backend/internal/id"
)

var ErrUnsupportedContent = errors.New("unsupported content type")

func getIDs(r *http.Request, keys ...string) (project id.ID, ids []int, err error) {
	vars := mux.Vars(r)

	v, ok := vars[projectID]
	if !ok {
		return "", nil, errors.New("project ID not provided")
	}
	prID, err := id.Parse(v)
	if err != nil {
		return "", nil, err
	}

	ids = make([]int, len(keys))
	for i, k := range keys {
		ids[i], err = strconv.Atoi(vars[k])
		if err != nil {
			return "", nil, err
		}
	}
	return prID, ids, nil
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
	return fmt.Errorf("%w: expected "+contentJSON, ErrUnsupportedContent)
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
