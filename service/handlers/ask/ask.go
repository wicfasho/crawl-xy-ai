package ask

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wicfasho/crawl-xy/scrapper"
)

type AskRequest struct {
	Question  string `json:"question" binding:"required,max=200"`
	SessionID string `json:"session_id" binding:"required,max=20"`
}

func Send(ctx *gin.Context) {
	var req AskRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	answer, err := scrapper.Ask(req.Question, req.SessionID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": *answer,
	})
}
