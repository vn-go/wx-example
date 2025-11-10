package core

import (
	"crypto/rand"
	"encoding/base64"
)

type secretService struct {
}

func (t *secretService) GenerateMasterKey() (string, error) {
	const size = 32 // 32 bytes = 256-bit key
	b := make([]byte, size)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func newSecretService() *secretService {
	return &secretService{}
}
