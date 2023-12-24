package Helpers

import "golang.org/x/crypto/bcrypt"

func CreateHashFromPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
