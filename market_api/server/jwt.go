package server

import (
	jwt "github.com/dgrijalva/jwt-go"
	//	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	SIGNING_KEY string = "abcdefggfedcbaxyz"
)

type AuthClaims struct {
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAuthToken(email, publicKey string) (string, error) {
	claims := AuthClaims{
		email,
		publicKey,
		jwt.StandardClaims{
			ExpiresAt: 0, // for now, just never expires todo
			Issuer:    "leopark",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SIGNING_KEY))
}

func ParseAuthToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNING_KEY), nil
	})
	claims, ok := token.Claims.(*AuthClaims)
	if ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("failed parse the token:" + err.Error())
}
