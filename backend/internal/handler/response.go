package handler

import (
	"event-scope/backend/internal/model"
	"time"

	"github.com/gin-gonic/gin"
)

func Res(ctx *gin.Context, status int, data any, err error) {
	msg := ""
	if err != nil {
		msg = err.Error()
	}

	ctx.JSON(status, model.APIResponse{
		Data:      data,
		Error:     msg,
		Timestamp: time.Now().In(time.FixedZone("Asia/Taipei", 8*60*60)),
	})
}
