package service

import (
	"gopastebin/consts"
	"log"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) string {
	salt, time, memory, threads, keyLen := consts.GetArgonOptions()
	hashedPassword := argon2.Key([]byte(password), salt, time, memory, threads, keyLen)

	return string(hashedPassword)
}


func VerifyPassword(password string, hashedPassword string) bool {

	log.Println("password: ", password)
	log.Println("hashedPassword: ", hashedPassword)
	return HashPassword(password) == hashedPassword
}