package pwd

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"core/services/cacher"
)

type bcryptPasswordService struct {
	cache cacher.CacheService
}

//var bcryptPassword = &bcryptPasswordService{}

// NewBcryptPasswordService tạo và trả về một instance mới của bcryptPasswordService.
// func (p *passwordServiceType) Bcrypt() passwordService {

// 	return bcryptPassword

// }

// HashPassword tạo một hash bcrypt an toàn cho mật khẩu đã cho.
// Chi phí (cost) mặc định của bcrypt là 10, bạn có thể điều chỉnh nó.
func (s *bcryptPasswordService) HashPassword(username, pass string) (string, error) {
	// bcrypt.DefaultCost là 10.
	key := strings.ToLower(username) + "@" + pass

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(key), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

// ComparePassword kiểm tra xem một mật khẩu dạng plain-text có khớp với hash bcrypt đã cho hay không.
func (s *bcryptPasswordService) ComparePassword(ctx context.Context, tenant, username, pass string, hashPass string) (bool, error) {
	cacheItem := &ComparePasswordCacheItem{}
	if err := s.cache.GetObject(ctx, tenant, username, &cacheItem); err == nil {
		return cacheItem.Ok, nil
	}
	key := strings.ToLower(username) + "@" + pass
	ret := false

	// So sánh mật khẩu plain-text với hash. Bcrypt sẽ tự động xử lý salt và các tham số khác.

	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(key))
	if err != nil {
		// Nếu lỗi là ErrMismatchedHashAndPassword, tức là mật khẩu không khớp.
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		// Các lỗi khác là lỗi thực sự.
		return false, fmt.Errorf("failed to compare password and hash: %w", err)
	}
	if err := s.cache.AddObject(ctx, tenant, username, &ComparePasswordCacheItem{
		Tenant:   tenant,
		Username: strings.ToLower(username),
		Ok:       true,
	}, 4); err != nil {
		return ret, err
	}

	return true, nil
}

type noEncryptPass struct {
}

func (n *noEncryptPass) HashPassword(pass string) (string, error) {
	// bcrypt.DefaultCost là 10.
	return pass, nil
}
func (n *noEncryptPass) ComparePassword(pass string, hashPass string) (bool, error) {
	// bcrypt.DefaultCost là 10.
	return pass == hashPass, nil
}
