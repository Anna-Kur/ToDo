package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
				return
			}

			if !validToken(cookie.Value) {
				http.Error(w, `{"error":"Authentication required"}`, http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}

func validToken(tokenString string) bool {
	key := []byte("secret")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	currentPassword := os.Getenv("TODO_PASSWORD")
	currentHash := sha256.Sum256([]byte(currentPassword))
	currentHashStr := hex.EncodeToString(currentHash[:])

	tokenHash, ok := claims["password_hash"].(string)
	if !ok || tokenHash != currentHashStr {
		return false
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return false
	}

	return true
}
