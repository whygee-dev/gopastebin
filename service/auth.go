package service

import (
	"gopastebin/consts"

	"github.com/golang-jwt/jwt"
)

func CreateToken() (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{}).SignedString(consts.GetSecret())

	if err != nil {
		return "", err
	}

	return token, nil
}