package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPass(t *testing.T) {
	argon := NewArgon2PasswordService()
	hash, err := argon.HashPassword("qwert45677#$%%")
	assert.NoError(t, err)
	ok, err := argon.ComparePassword("qwert45677#$%%", hash)
	assert.NoError(t, err)
	assert.True(t, ok)

	bcrPass := NewBcryptPasswordService()
	hash, err = bcrPass.HashPassword("qwert45677#$%%")
	assert.NoError(t, err)
	ok, err = bcrPass.ComparePassword("qwert45677#$%%", hash)
	assert.NoError(t, err)
	assert.True(t, ok)

}
