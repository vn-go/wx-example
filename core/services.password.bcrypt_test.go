package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vn-go/bx"
)

func TestBcryptPasswordService(t *testing.T) {
	argon := PasswordService.Bcrypt
	hash, err := argon.HashPassword("qwert45677#$%%")
	assert.NoError(t, err)
	ok, err := argon.ComparePassword("qwert45677#$%%", hash)
	assert.NoError(t, err)
	assert.True(t, ok)

}
func BenchmarkBcryptPasswordService(t *testing.B) {
	t.Run("call function", func(t *testing.B) {
		for i := 0; i < t.N; i++ {

			hash, _ := PasswordService.Bcrypt.HashPassword("qwert45677#$%%")

			PasswordService.Bcrypt.ComparePassword("qwert45677#$%%", hash)

		}
	})
	t.Run("no call function", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			svc := PasswordService.Bcrypt
			hash, _ := svc.HashPassword("qwert45677#$%%")

			svc.ComparePassword("qwert45677#$%%", hash)

		}
	})
}
func BenchmarkOnceCall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bx.OnceCall[string]("test", func() (string, error) {
			return "hello", nil
		})
	}
}

func BenchmarkDirect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = "hello"
	}
}
