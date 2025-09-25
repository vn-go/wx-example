package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgon2PasswordService(t *testing.T) {
	argon := &argon2PasswordService{}
	hash, err := argon.HashPassword("admin", "qwert45677#$%%")
	assert.NoError(t, err)
	ok, err := argon.ComparePassword(t.Context(), "hrm", "admin", "qwert45677#$%%", hash)
	assert.NoError(t, err)
	assert.True(t, ok)
	bcrypt := &bcryptPasswordService{}
	hash, err = bcrypt.HashPassword("root", "root")
	ok, err = bcrypt.ComparePassword(t.Context(), "hrm", "root", "root", hash)
	assert.NoError(t, err)
	assert.True(t, ok)

}

func BenchmarkCompare(t *testing.B) {
	t.Run("Bcrypt", func(t *testing.B) {
		Bcrypt := &bcryptPasswordService{}
		for i := 0; i < t.N; i++ {

			hash, _ := Bcrypt.HashPassword("admin", "qwert45677#$%%")

			Bcrypt.ComparePassword(t.Context(), "hrm", "admin", "qwert45677#$%%", hash)

		}
	})
	t.Run("Argon2", func(t *testing.B) {
		Argon2 := &argon2PasswordService{}
		for i := 0; i < t.N; i++ {

			hash, _ := Argon2.HashPassword("admin", "qwert45677#$%%")

			Argon2.ComparePassword(t.Context(), "hrm", "admin", "qwert45677#$%%", hash)

		}
	})
}
