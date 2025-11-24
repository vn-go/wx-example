package pwd

import (
	"context"
	"core/services/cacher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type argon2PasswordService struct {
	cache cacher.CacheService
}

// NewArgon2PasswordService creates and returns a new instance of argon2PasswordService.
// func (p *passwordServiceType) Argon2() passwordService {
// 	ret, _ := bx.OnceCall[passwordServiceType]("NewArgon2PasswordService", func() (passwordService, error) {
// 		return &argon2PasswordService{}, nil
// 	})
// 	return ret

// }

// HashPassword generates a secure Argon2 hash for a given password.
// The resulting string is a full hash, including salt and parameters,
// encoded in a standard format.
func (s *argon2PasswordService) HashPassword(username, pass string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Parameters for Argon2id. You can tune these for your security needs.
	var (
		memory      = uint32(64 * 1024) // 64 MB
		iterations  = uint32(3)
		parallelism = uint8(2)
		keyLen      = uint32(32) // 32 bytes for the key
	)
	pwdHas := fmt.Sprintf("%s@%s", strings.ToLower(username), pass)
	hash := argon2.IDKey([]byte(pwdHas), salt, iterations, memory, parallelism, keyLen)

	// Format the hash string for storage.
	// This format is a common Argon2 string representation.
	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		memory,
		iterations,
		parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	)

	return encodedHash, nil
}

type ComparePasswordCacheItem struct {
	Tenant   string
	Username string
	Ok       bool
}

// ComparePassword checks if a plain-text password matches a given Argon2 hash.
// It parses the hash string to extract the salt and parameters, then
// re-computes the hash and performs a constant-time comparison.
func (s *argon2PasswordService) ComparePassword(ctx context.Context, tenant, username, pass, hashPass string) (bool, error) {

	ret := false
	cacheItem := &ComparePasswordCacheItem{}
	if err := s.cache.GetObject(ctx, tenant, username, &cacheItem); err == nil {
		return cacheItem.Ok, nil
	}
	parts := strings.Split(hashPass, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash string format: expected 6 parts, got %d", len(parts))
	}

	var (
		//version     int
		memory      uint32
		iterations  uint32
		parallelism uint8
	)

	// Parse parameters from the second part
	// Example: m=65536,t=3,p=2
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return false, fmt.Errorf("invalid parameter format")
	}

	_, err := fmt.Sscanf(params[0], "m=%d", &memory)
	if err != nil {
		return false, fmt.Errorf("failed to parse memory param: %w", err)
	}
	_, err = fmt.Sscanf(params[1], "t=%d", &iterations)
	if err != nil {
		return false, fmt.Errorf("failed to parse iterations param: %w", err)
	}
	_, err = fmt.Sscanf(params[2], "p=%d", &parallelism)
	if err != nil {
		return false, fmt.Errorf("failed to parse parallelism param: %w", err)
	}

	// Decode salt and hash
	decodedSalt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, fmt.Errorf("failed to decode salt: %w", err)
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, fmt.Errorf("failed to decode hash: %w", err)
	}
	pwdHas := fmt.Sprintf("%s@%s", strings.ToLower(username), pass)
	// Re-compute the hash with the extracted parameters.
	recomputedHash := argon2.IDKey(
		[]byte(pwdHas),
		decodedSalt,
		iterations,
		memory,
		parallelism,
		uint32(len(decodedHash)),
	)

	// Use a constant-time comparison to prevent timing attacks.
	if subtle.ConstantTimeCompare(recomputedHash, decodedHash) == 1 {
		return true, nil
	}
	if err := s.cache.AddObject(ctx, tenant, username, &ComparePasswordCacheItem{
		Tenant:   tenant,
		Username: strings.ToLower(username),
		Ok:       true,
	}, 4); err != nil {
		return ret, nil
	}

	return false, nil
}
