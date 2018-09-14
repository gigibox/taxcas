package app

import (
	"github.com/gin-gonic/gin"

	"taxcas/pkg/e"
)

type Gin struct {
	C *gin.Context
}

type ResponseMsg struct {
	success bool
	msg     string
	data    interface{}
}

func (g *Gin) Response(httpCode int, result bool, errCode int, data interface{}) {
	g.C.JSON(httpCode, gin.H{
		"success": result,
		"msg":  e.GetMsg(errCode),
		"data": data,
	})

	return
}
