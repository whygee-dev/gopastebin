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

func DecodeToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return consts.GetSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}