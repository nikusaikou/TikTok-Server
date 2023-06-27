package middleware

import (
	"TikTokServer/pkg/auth"
	"TikTokServer/pkg/errorcode"
	response "TikTokServer/pkg/respponse"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRequest := c.PostForm("token")
		if tokenRequest == "" {
			tokenRequest = c.Query("token")
		}

		userID, err := auth.GetUserIDByToken(tokenRequest)
		if err != nil && userID == int64(-1) {
			response.Fail(c, errorcode.ErrHttpTokenInvalid, nil)
			c.Abort()
		}
		c.Set("userID", userID)
		c.Next()
	}
}
