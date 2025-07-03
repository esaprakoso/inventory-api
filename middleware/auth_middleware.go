package middleware

import (
	"net/http"
	"pos/config"
	"pos/handlers"
	"strings"

	"github.com/gin-gonic/gin"
	jwtGo "github.com/golang-jwt/jwt/v4"
)

func Protected() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or malformed JWT"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Missing or malformed JWT"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims := &handlers.AuthClaims{}
		token, err := jwtGo.ParseWithClaims(tokenString, claims, func(token *jwtGo.Token) (any, error) {
			return []byte(config.LoadConfig("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired JWT"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.Issuer)
		c.Next()
	}
}
