package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserId
func GetUserId(ctx *gin.Context) (userId int) {
	// handle user Id
	userIdRaw, exist := ctx.Get("user_id")

	if !exist {
		userId = 0 // no auth user
	} else {
		switch v := userIdRaw.(type) {
		case int:
			userId = v
		case int64:
			userId = int(v)
		case float64:
			userId = int(v)
		case string:
			fmt.Sscanf(v, "%d", &userId)
		default:
			userId = 0
		}
	}
	return
}
