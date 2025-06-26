package middlewares

import (
	"strings"
	"time"

	sharedUtil "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// security headers
func SecurityHeaders() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Cache-Control", "no-cache, no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}

}

// cors header set
func CorsHeader() gin.HandlerFunc {

	// string to []string
	origins := strings.Split(sharedUtil.GetEnv("ALLOWED_ORIGINS", ""), ",")

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cache-Control", "X-Requested-With", "User-Agent", " Content-Length", "Accept-Encoding", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
