package repository

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

func Compare(hashed, plane string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plane))
	return err == nil
}
