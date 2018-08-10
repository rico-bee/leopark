package server

import (
	crypto "github.com/rico-bee/leopark/crypto"
	"log"
	"net/http"
)

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		log.Println("auth middleware:" + tokenString)
		if len(tokenString) == 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing Authorization Header"))
			return
		}
		auth, err := crypto.ParseAuthToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Error verifying JWT token: " + err.Error()))
			return
		}
		log.Println("auth middleware done:" + auth.PrivateKey)
		r.Header.Set("email", auth.Email)
		r.Header.Set("publicKey", auth.PublicKey)
		next.ServeHTTP(w, r)
	})
}
