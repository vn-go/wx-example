package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	cfg, err := Config.Load("config.yaml")
	assert.NoError(t, err)
	assert.NotEmpty(t, cfg)
}
func BenchmarkLoadConfig(t *testing.B) {
	for i := 0; i < t.N; i++ {
		Config.Load("config.yaml")

	}

}
