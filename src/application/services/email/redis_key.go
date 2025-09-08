package email

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	EmailTokenKeyPrefix = "email_token:%s"
)

var (
	UserTokenExpireDuration = time.Hour * 1
)

func GetEmailTokenKey(email string) string {
	return fmt.Sprintf(EmailTokenKeyPrefix, email)
}
func GetEnvAsInt64OrDefault(key string, defaultValue int64) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(intValue)
		}
	}
	return time.Duration(defaultValue)
}
