package core

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type bcryptPasswordService struct{}

//var bcryptPassword = &bcryptPasswordService{}

// NewBcryptPasswordService tạo và trả về một instance mới của bcryptPasswordService.
// func (p *passwordServiceType) Bcrypt() passwordService {

// 	return bcryptPassword

// }

// HashPassword tạo một hash bcrypt an toàn cho mật khẩu đã cho.
// Chi phí (cost) mặc định của bcrypt là 10, bạn có thể điều chỉnh nó.
func (s *bcryptPasswordService) HashPassword(pass string) (string, error) {
	// bcrypt.DefaultCost là 10.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

// ComparePassword kiểm tra xem một mật khẩu dạng plain-text có khớp với hash bcrypt đã cho hay không.
func (s *bcryptPasswordService) ComparePassword(pass string, hashPass string) (bool, error) {
	// So sánh mật khẩu plain-text với hash. Bcrypt sẽ tự động xử lý salt và các tham số khác.
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(pass))
	if err != nil {
		// Nếu lỗi là ErrMismatchedHashAndPassword, tức là mật khẩu không khớp.
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		// Các lỗi khác là lỗi thực sự.
		return false, fmt.Errorf("failed to compare password and hash: %w", err)
	}

	return true, nil
}
