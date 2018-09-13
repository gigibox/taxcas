package app

import (
	"github.com/gin-gonic/gin"

	"taxcas/pkg/e"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) Response(httpCode int, result bool, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"success": result,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})

	return
}
