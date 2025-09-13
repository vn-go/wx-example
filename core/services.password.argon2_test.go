package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2PasswordService(t *testing.T) {
	argon := PasswordService.Argon2
	hash, err := argon.HashPassword("qwert45677#$%%")
	assert.NoError(t, err)
	ok, err := argon.ComparePassword("qwert45677#$%%", hash)
	assert.NoError(t, err)
	assert.True(t, ok)

}
func BenchmarkArgon2PasswordService(t *testing.B) {
	for i := 0; i < t.N; i++ {
		argon := PasswordService.Argon2
		hash, _ := argon.HashPassword("qwert45677#$%%")

		argon.ComparePassword("qwert45677#$%%", hash)

	}

}
func BenchmarkCompare(t *testing.B) {
	t.Run("Bcrypt", func(t *testing.B) {
		for i := 0; i < t.N; i++ {

			hash, _ := PasswordService.Bcrypt.HashPassword("qwert45677#$%%")

			PasswordService.Bcrypt.ComparePassword("qwert45677#$%%", hash)

		}
	})
	t.Run("Argon2", func(t *testing.B) {
		for i := 0; i < t.N; i++ {

			hash, _ := PasswordService.Argon2.HashPassword("qwert45677#$%%")

			PasswordService.Argon2.ComparePassword("qwert45677#$%%", hash)

		}
	})
}
