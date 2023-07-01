package middleware

import (
	"TikTokServer/pkg/auth"
	"TikTokServer/pkg/errorcode"
	response "TikTokServer/pkg/response"
	"TikTokServer/pkg/tlog"

	"github.com/gin-gonic/gin"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenRequest := c.PostForm("token")
		if tokenRequest == "" {
			tokenRequest = c.Query("token")
		}

		userID, err := auth.GetUserIDByToken(tokenRequest)

		tlog.Debugf("middleware userID: %v", userID)
		tlog.Debugf("middleware err: %v", err)
		if err != nil || userID == int64(-1) {
			errCode := errorcode.ErrHttpTokenInvalid
			errCode.SetError(err)
			response.Fail(c, errCode, nil)
			c.Abort()
		}
		c.Set("userID", userID)
		c.Next()
	}
}
