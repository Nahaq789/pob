package repository

import (
	"crypto/sha256"

	"golang.org/x/crypto/bcrypt"
)

func Hash(v string) (string, error) {
	hash := sha256.Sum256([]byte(v))
	bcryptHashed, err := bcrypt.GenerateFromPassword(hash[:], bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bcryptHashed), nil
}

func Compare(hashed, plane string) bool {
	hash := sha256.Sum256([]byte(plane))
	err := bcrypt.CompareHashAndPassword([]byte(hashed), hash[:])
	return err == nil
}
