package rpc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	//	"github.com/hyperledger/sawtooth-sdk-go/signing"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
)

const (
	AES_KEY     string = "ffffffffffffffffffffffffffffffff"
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

func EncryptKey(aesKey, publicKey, privateKey string) ([]byte, error) {
	keyBytes := []byte(privateKey)
	if len(keyBytes)%aes.BlockSize != 0 {
		return nil, errors.New("key must be a multiple of block size")
	}
	aesKeyBytes, _ := hex.DecodeString(aesKey)
	block, err := aes.NewCipher(aesKeyBytes)
	if err != nil {
		return nil, err
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(keyBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(ciphertext[aes.BlockSize:], keyBytes)
	return ciphertext, nil
}

func DecryptKey(aesKey, publicKey, privateKey string) (string, error) {
	aesKeyBytes, _ := hex.DecodeString(aesKey)
	cipherBytes, _ := hex.DecodeString(privateKey)
	block, err := aes.NewCipher(aesKeyBytes)
	if err != nil {
		return "", err
	}
	if len(cipherBytes) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := cipherBytes[:aes.BlockSize]
	cipherBytes = cipherBytes[aes.BlockSize:]
	// CBC mode always works in whole blocks.
	if len(cipherBytes)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(cipherBytes, cipherBytes)
	return hex.EncodeToString(cipherBytes), nil
}
