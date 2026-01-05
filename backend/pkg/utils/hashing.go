package utils

import "golang.org/x/crypto/bcrypt"

func HashingValue(password string) (string, error) {
	hashValue, err :=bcrypt.GenerateFromPassword([]byte (password), 13)
	return string(hashValue), err
}

func ChcekPassword(hashValue, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashValue), []byte(password))

	return err==nil
}