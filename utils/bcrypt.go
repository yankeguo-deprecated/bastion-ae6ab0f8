package utils

import "golang.org/x/crypto/bcrypt"

func BcryptGenerate(password string) (string, error) {
	buf, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(buf), err
}

func BcryptValidate(hashedPassword string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
