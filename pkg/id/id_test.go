package id

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID(t *testing.T) {
	id, err := New()
	require.NoError(t, err)
	assert.Len(t, id, length)
	for _, r := range id {
		if (r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9') {
			continue
		}
		assert.Fail(t, fmt.Sprintf("unexpected character: %c", r), id)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		id  string
		err error
	}{
		{"", ErrInvalidID},
		{"1234", ErrInvalidID},
		{"asd31234", nil},
		{"as__d234", ErrInvalidID},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			_, err := Parse(test.id)
			assert.ErrorIs(t, err, test.err)
		})
	}
}
