package model

import (
	"errors"
	"unicode"
)

type scanner interface {
	Scan(...interface{}) error
}

var ErrNotFound = errors.New("not found")

func stringMadeOf(s string, ranges ...*unicode.RangeTable) bool {
	for _, r := range s {
		if !unicode.In(r, ranges...) {
			return false
		}
	}
	return true
}
