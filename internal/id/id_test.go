package id

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
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
