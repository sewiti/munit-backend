package id

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	length  = 8
	symbols = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
)

var lenSymbols = big.NewInt(int64(len(symbols)))

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
	if len(s) != length {
		return "", fmt.Errorf("invalid ID length: %s", s)
	}
	return ID(s), nil
}
