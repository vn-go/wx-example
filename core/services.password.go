package core

type passwordService interface {
	HashPassword(pass string) (string, error)
	ComparePassword(pass string, hashPass string) (bool, error)
}
type passwordServiceType struct {
	Bcrypt passwordService
	Argon2 passwordService
}

var PasswordService = &passwordServiceType{
	Bcrypt: &bcryptPasswordService{},
	Argon2: &argon2PasswordService{},
}
