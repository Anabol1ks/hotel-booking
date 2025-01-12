package auth

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется авторизация"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			// Check if it's a system token
			systemToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("SYSTEM_TOKEN_SECRET")), nil
			})

			if err != nil || !systemToken.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный токен"})
				c.Abort()
				return
			}

			claims := systemToken.Claims.(jwt.MapClaims)
			if claims["system"] != true {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный системный токен"})
				c.Abort()
				return
			}

			c.Set("system", true)
			c.Next()
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := uint(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Set("role", claims["role"])

		c.Next()
	}
}
