package utils

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func IsStrongPassword(pw string) bool {
	var hasUpper, hasLower, hasDigit, hasSpecial bool

	for _, char := range pw {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	lengthOK := len(pw) >= 8
	return hasUpper && hasLower && hasDigit && hasSpecial && lengthOK
}
