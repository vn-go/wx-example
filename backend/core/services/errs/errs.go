package errs

type ErrService struct {
}

func (e *ErrService) IsForbidden(err error) bool {
	if e, ok := err.(*Err); ok {
		return e.Typ == ErrForbidden
	}
	return false
}

func (e ErrService) Unauthenticate() error {
	return &Err{
		Message: "Unauthenticate",
		Typ:     ErrUnautheticate,
	}
}

func (e ErrService) ForbiddenErr() error {

	return &Err{
		Message: "forbidden",
		Typ:     ErrForbidden,
	}
}

func (e ErrService) BadRequest(msg string, args ...any) error {
	return &Err{
		Message: msg,
		Args:    args,
		Typ:     ErrBadRequest,
	}

}

func NewErrService() *ErrService {
	return &ErrService{}
}
