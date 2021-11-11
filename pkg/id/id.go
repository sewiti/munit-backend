package id

import (
	"crypto/rand"
	"errors"
	"math/big"
	"strings"
)

const (
	length  = 8
	symbols = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
)

var lenSymbols = big.NewInt(int64(len(symbols)))

var ErrInvalidID = errors.New("invalid id")

type ID string

func New() (ID, error) {
	id := make([]byte, length)
	for i := range id {
		n, err := rand.Int(rand.Reader, lenSymbols)
		if err != nil {
			return "", err
		}
		id[i] = symbols[n.Int64()]
	}
	return ID(id), nil
}

func Parse(s string) (ID, error) {
	if err := ID(s).Validate(); err != nil {
		return "", err
	}
	return ID(s), nil
}

// Validate checks ID's validity.
// Returns ErrInvalidID if invalid.
func (id ID) Validate() error {
	if len(id) != length {
		return ErrInvalidID
	}
	for _, r := range id {
		if !strings.ContainsRune(symbols, r) {
			return ErrInvalidID
		}
	}
	return nil
}
