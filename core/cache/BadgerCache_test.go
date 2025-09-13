package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBagger(t *testing.T) {
	bg, err := NewBadgerCache("Badger")
	assert.NoError(t, err)
	type testData struct {
		A string
		B int
	}
	data := testData{
		A: "dsadsadasd",
		B: 12343,
	}
	bg.Set(t.Context(), "test-001", &data, time.Hour)
}
