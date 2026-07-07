package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	expectedPassword := os.Getenv("TODO_PASSWORD")
	if expectedPassword == "" {
		writeJSONError(w, "Authentication not configured", http.StatusInternalServerError)
		return
	}

	if req.Password != expectedPassword {
		writeJSONError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(req.Password)
	if err != nil {
		writeJSONError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token": token,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func generateToken(password string) (string, error) {
	key := []byte("secret")

	hash := sha256.Sum256([]byte(password))
	passwordHash := hex.EncodeToString(hash[:])

	claims := jwt.MapClaims{
		"password_hash": passwordHash,
		"exp":           time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("Failed to sign jwt: %w", err)

	}

	return signedToken, nil
}
