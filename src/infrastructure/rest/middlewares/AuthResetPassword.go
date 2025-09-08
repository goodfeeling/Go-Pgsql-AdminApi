package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// AuthResetPassword 使用重置密码认证中间件
func AuthResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}

		resetSecret := os.Getenv("JWT_RESET_SECRET")
		if resetSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT_RESET_SECRET not configured"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(resetSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < jwt.TimeFunc().Unix() {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		if t, ok := claims["type"].(string); ok {
			if t != "reset" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Token type mismatch"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Missing token type"})
			c.Abort()
			return
		}

		c.Next()
	}
}
