package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
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
// @router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}
	valid := validation.Validation{}

	username := c.Query("username")
	password := c.Query("password")

	authService := auth_service.Auth{Username: username, Password: password}
	ok, _ := valid.Valid(authService)

	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	isExist, err := authService.Check()
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusOK, false, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	appG.C.SetCookie("token", token, 60*30, "/", "", false, true)
	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"token": token,
	})
}
