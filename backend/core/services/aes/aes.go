package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Dịch vụ để đóng gói các hàm mã hóa
type AesService struct{}

// Độ dài khóa AES-256 (32 byte)
const KeyLength = 32

// --- 1. Hàm Sinh Khóa Ngẫu nhiên (Generate Random Key) ---
// Khóa ngẫu nhiên phải có độ dài 32 byte cho AES-256
func (s *AesService) GenerateRandomKey() (string, error) {
	key := make([]byte, KeyLength)
	// Đọc 32 byte ngẫu nhiên từ nguồn ngẫu nhiên an toàn về mặt mật mã
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	// Trả về khóa dưới dạng Base64 URL-safe string
	return base64.URLEncoding.EncodeToString(key), nil
}

// --- 2. Hàm Mã hóa (Encrypt) ---
// Mã hóa dữ liệu bằng khóa và chế độ GCM.
// Khóa (key) phải là chuỗi Base64 URL-safe.
// Đầu ra là chuỗi Base64 URL-safe chứa Nonce (IV) + Dữ liệu đã mã hóa + Tag GCM
func (s *AesService) Encrypt(data string, key string) (string, error) {
	// 1. Giải mã khóa từ Base64
	keyBytes, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New("invalid key format")
	}
	if len(keyBytes) != KeyLength {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}

	// 2. Tạo khối mã hóa AES
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	// 3. Tạo GCM (Galois/Counter Mode)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 4. Tạo Nonce (IV - Initialization Vector) ngẫu nhiên
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 5. Thực hiện mã hóa (Nonce được gắn vào đầu đầu ra đã mã hóa)
	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	// 6. Trả về dưới dạng Base64
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// --- 3. Hàm Giải mã (Decrypt) ---
// Giải mã dữ liệu đã mã hóa bằng khóa bí mật.
func (s *AesService) Decrypt(data string, key string) (string, error) {
	// 1. Giải mã khóa từ Base64
	keyBytes, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return "", errors.New("invalid key format")
	}
	if len(keyBytes) != KeyLength {
		return "", errors.New("key length must be 32 bytes for AES-256")
	}

	// 2. Giải mã dữ liệu đầu vào Base64
	ciphertext, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("invalid ciphertext format")
	}

	// 3. Tạo khối AES và GCM
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 4. Tách Nonce
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	// Nonce được lưu ở đầu chuỗi mã hóa
	nonce, encryptedDataWithTag := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// 5. Thực hiện giải mã và xác thực (Open)
	// Nếu dữ liệu bị thay đổi, gcm.Open sẽ trả về lỗi
	plaintextBytes, err := gcm.Open(nil, nonce, encryptedDataWithTag, nil)
	if err != nil {
		return "", errors.New("decryption failed or data has been tampered with (GCM Tag invalid)")
	}

	// 6. Trả về dữ liệu gốc
	return string(plaintextBytes), nil
}

func NewAesService() *AesService {
	return &AesService{}
}
