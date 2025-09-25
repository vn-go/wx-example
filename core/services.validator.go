package core

import (
	"errors"
	"fmt"
	"unicode"
)

type validatorService struct {
}

// Validation rules (tùy chỉnh nếu cần)
const (
	MinUsernameLen = 3
	MaxUsernameLen = 64
)

// Err definitions
var (
	ErrEmptyUsername       = errors.New("username is empty")
	ErrTooShort            = fmt.Errorf("username must be at least %d characters", MinUsernameLen)
	ErrTooLong             = fmt.Errorf("username must be at most %d characters", MaxUsernameLen)
	ErrInvalidChar         = errors.New("username contains invalid character(s)")
	ErrConsecutiveDots     = errors.New("username must not contain consecutive dots ('..')")
	ErrStartsOrEndsWithDot = errors.New("username must not start or end with a dot ('.')")
)

// ValidateUsername checks whether a username meets the policy.
// Returns nil if valid, otherwise a descriptive error.
func (v *validatorService) ValidateUsername(u string) error {
	if len(u) == 0 {
		return ErrEmptyUsername
	}
	if len(u) < MinUsernameLen {
		return ErrTooShort
	}
	if len(u) > MaxUsernameLen {
		return ErrTooLong
	}

	var prev rune
	for i, r := range u {
		// Allowed:
		//  - Unicode letters (categories L*)
		//  - ASCII digits 0-9
		//  - underscore '_', hyphen '-', dot '.'
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.' {
			// ok
		} else {
			return ErrInvalidChar
		}

		// No consecutive dots
		if r == '.' && prev == '.' {
			return ErrConsecutiveDots
		}
		prev = r

		// check first/last char for dot later; we can also check here for index
		if i == 0 && r == '.' {
			return ErrStartsOrEndsWithDot
		}
	}

	// check last character not dot
	if last := rune(u[len(u)-1]); last == '.' {
		return ErrStartsOrEndsWithDot
	}

	return nil
}
func newValidatorService() *validatorService {
	return &validatorService{}
}
