package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MiddlewareProvider struct {
	RedisClient *redis.Client
	DB          *gorm.DB
}

func NewMiddlewareProvider(redisClient *redis.Client, db *gorm.DB) *MiddlewareProvider {
	return &MiddlewareProvider{
		RedisClient: redisClient,
		DB:          db,
	}
}

func (mp *MiddlewareProvider) AuthJWTMiddleware() gin.HandlerFunc {
	return AuthJWTMiddlewareWithRedis(mp.RedisClient, mp.DB)
}

func (mp *MiddlewareProvider) OptionalAuthMiddleware() gin.HandlerFunc {
	return OptionalAuthMiddlewareWithRedis(mp.RedisClient, mp.DB)
}

func (mp *MiddlewareProvider) UrlAuthMiddleware() gin.HandlerFunc {
	return UrlAuthMiddlewareWithRedis(mp.RedisClient, mp.DB)
}
