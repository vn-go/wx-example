package errs

import "fmt"

type ErrType int

const (
	ErrSytem ErrType = iota
	ErrBadRequest
	ErrForbidden
	ErrUnautheticate
)

func (t ErrType) String() string {
	switch t {
	case ErrSytem:
		return "System"
	case ErrBadRequest:
		return "Bad Request"
	case ErrForbidden:
		return "Forbidden"
	case ErrUnautheticate:
		return "ErrUnautheticate"
	default:
		return "Unknown"
	}
}

type Err struct {
	Message string
	Args    []any
	Typ     ErrType
}

// Error implements error.
func (e *Err) Error() string {
	args := append(e.Args, e.Typ.String())
	return fmt.Sprintf(e.Message, args...) + "\n" + e.Typ.String()
}
