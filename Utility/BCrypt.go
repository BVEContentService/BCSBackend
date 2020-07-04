package Utility

import (
    "golang.org/x/crypto/bcrypt"
)

func BCryptCalculateHash(pass string) []byte {
    hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
    if err != nil { panic(err) }
    return hash
}

func BCryptValidateHash(pass string, hash []byte) bool {
    if len(hash)==0 && pass=="" { return true }
    err := bcrypt.CompareHashAndPassword(hash, []byte(pass))
    return err == nil
}