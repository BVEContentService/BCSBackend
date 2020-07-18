package Utility

import (
	"golang.org/x/crypto/bcrypt"
)

func BCryptCalculateHash(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func BCryptValidateHash(pass string, hash string) bool {
	if hash == "" && pass == "" {
		return true
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	return err == nil
}
