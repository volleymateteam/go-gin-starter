package utils

import (
	"math/rand"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
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

// New Random Password Generator
const passwordCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomPassword() string {
	length := 12
	password := make([]byte, length)
	for i := range password {
		password[i] = passwordCharset[seededRand.Intn(len(passwordCharset))]
	}
	return string(password)
}
