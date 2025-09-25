package broker

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	salt := make([]byte, 16)
	_, _ = rand.Read(salt)

	// memory = 64 MB, iterations = 3, parallelism = 4
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	return base64.RawStdEncoding.EncodeToString(hash)
}
func BenchmarkArgon2PasswordService(t *testing.B) {
	password := "a_very_secure_password_with_some_extra_length"
	for i := 0; i < t.N; i++ {
		_ = HashPassword(password)
		//assert.NoError(t, err)
		// argon := &argon2PasswordService{}
		// hash, _ := argon.HashPassword("admin", "qwert45677#$%%")

		// argon.ComparePassword(t.Context(), "hrm", "admin", "qwert45677#$%%", hash)

	}

}
