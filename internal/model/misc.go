package model

import (
	"errors"
	"unicode"

	"github.com/go-sql-driver/mysql"
)

type scanner interface {
	Scan(...interface{}) error
}

var ErrNotFound = errors.New("resource not found")

func stringMadeOf(s string, ranges ...*unicode.RangeTable) bool {
	for _, r := range s {
		if !unicode.In(r, ranges...) {
			return false
		}
	}
	return true
}

func isDuplicate(err error) bool {
	sqlErr, ok := err.(*mysql.MySQLError)
	return ok && sqlErr.Number == 1062
}
