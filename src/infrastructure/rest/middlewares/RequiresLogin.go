package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthJWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}

		accessSecret := os.Getenv("JWT_ACCESS_SECRET")
		if accessSecret == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT_ACCESS_SECRET not configured"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims := jwt.MapClaims{}
		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return []byte(accessSecret), nil
		})
		if err != nil {
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
			if t != "access" {
				c.JSON(http.StatusForbidden, gin.H{"error": "Token type mismatch"})
				c.Abort()
				return
			}
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "Missing token type"})
			c.Abort()
			return
		}

		if idFloat, ok := claims["id"].(float64); ok {
			id := int(idFloat)
			if id == 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "Missing or invalid user id"})
				c.Abort()
				return
			}
			c.Set("user_id", id)
		}
		if roleIdFloat64, ok := claims["role_id"].(float64); ok {
			id := int64(roleIdFloat64)
			if id == 0 {
				c.JSON(http.StatusForbidden, gin.H{"error": "Missing or invalid role id"})
				c.Abort()
				return
			}
			c.Set("role_id", id)
		}

		c.Next()
	}
}
