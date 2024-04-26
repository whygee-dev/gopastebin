package service

import (
	"gopastebin/consts"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	salt, time, memory, threads, keyLen := consts.GetArgonOptions()
	hashedPassword := argon2.Key([]byte(password), salt, time, memory, threads, keyLen)

	return string(hashedPassword)
}


func VerifyPassword(password string, hashedPassword string) bool {
	return HashPassword(password) == hashedPassword
}