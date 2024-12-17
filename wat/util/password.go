package wat

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hashedPassword, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(providedPassword))
	if err != nil {
		return false
	}
	return true
}
