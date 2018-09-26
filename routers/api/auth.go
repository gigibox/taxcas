package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/util"
	"taxcas/service/auth_service"
)

// @Summary 用户登陆
// @Tags 	认证授权
// @Param   username query string true "The username for login"
// @Param   password query string true "The password for login"
// @Success 200 {object} app.ResponseMsg "data":{"token":string}"
// @router  /api/auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}

	username := c.Query("username")
	password := c.Query("password")

	if len(username) > 30 || len(password) > 30 {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	isSuccess, code := auth_service.CheckAuth(username, password)
	if isSuccess {
		c.Header("Authorization", util.GenerateToken("admin", username))
	}

	appG.Response(http.StatusOK, isSuccess, code, nil)
	return
}
