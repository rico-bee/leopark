package crypto

import (
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type AuthClaims struct {
	Email     string `json:"email"`
	PublicKey string `json:"public_key"`
	jwt.StandardClaims
}

type AuthInfo struct {
	Email      string `gorethink:"email"`
	PublicKey  string `gorethink:"publicKey"`
	PwdHash    string `gorethink:"pwdHash,omitempty"`
	PrivateKey string `gorethink:"privateKey,omitempty"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateAuthToken(auth *AuthInfo) (string, error) {

	claims := AuthClaims{
		auth.Email,
		auth.PublicKey,
		jwt.StandardClaims{
			ExpiresAt: 0, // for now, just never expires todo
			Issuer:    "leopark",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SIGNING_KEY))
}

func ParseAuthToken(tokenString string) (*AuthInfo, error) {
	claims := &AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGNING_KEY), nil
	})
	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		authInfo := &AuthInfo{
			Email:     claims.Email,
			PublicKey: claims.PublicKey,
		}
		log.Println("email:" + authInfo.Email)
		return authInfo, nil
	} else {
		log.Println("failed parse the token:" + err.Error())
		return nil, err
	}
}
