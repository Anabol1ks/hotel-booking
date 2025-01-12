package auth

import (
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	systemToken string
	tokenMutex  sync.RWMutex
	tokenTTL    = 72 * time.Hour
)

func GetSystemToken() string {
	tokenMutex.RLock()
	currentToken := systemToken
	tokenMutex.RUnlock()

	if currentToken == "" {
		return generateSystemToken()
	}

	// Verify token validity
	token, _ := jwt.Parse(currentToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SYSTEM_TOKEN_SECRET")), nil
	})

	if token == nil || !token.Valid {
		return generateSystemToken()
	}

	return currentToken
}

func generateSystemToken() string {
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"system": true,
		"exp":    time.Now().Add(tokenTTL).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(os.Getenv("SYSTEM_TOKEN_SECRET")))
	systemToken = tokenString
	return tokenString
}
