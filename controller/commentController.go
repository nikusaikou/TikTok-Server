package controller

import (
	"TikTokServer/pkg/log"

	"github.com/gin-gonic/gin"
)

func CommentAction(c *gin.Context) {
	log.Info("CommentAction")
}

func GetCommentList(c *gin.Context) {
	log.Info("CommentList")
}