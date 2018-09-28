package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/util"
	"taxcas/service/auth_service"
)

// @Summary 后台登陆
// @Tags 	开放接口
// @Description 登陆成功后返回 "success":true, 并在header中返回token, "Authorization: token". 后续访问接口需要在header中添加该字段 "Authorization: Bearer token"
// @Param   username query string true "The username for login"
// @Param   password query string true "The password for login"
// @Success 200 {object} app.ResponseMsg
// @router  /api/admin/login [get]
func Login(c *gin.Context) {
	appG := app.Gin{C: c}

	username := c.Query("username")
	password := c.Query("password")

	if len(username) > 30 || len(password) > 30 {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, nil)
		return
	}

	isSuccess, code := auth_service.CheckAuth(username, password)
	if isSuccess {
		c.Header("Authorization", util.GenerateToken("admin", util.EncodeMD5(username)))
	}

	appG.Response(http.StatusOK, isSuccess, code, nil)
	return
}


type changepwd struct {
	UserName    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
// @Summary 修改密码
// @Tags 	后台管理
// @Security ApiKeyAuth
// @Produce  json
// @Param   data body admin.changepwd true "用户名 | 新旧密码"
// @Success 200 {object} app.ResponseMsg
// @router  /api/v1/admin/password [put]
func ChangePassword(c *gin.Context) {
	appG := app.Gin{C: c}

	params := changepwd{}
	err := c.BindJSON(&params)
	if err != nil {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, err)
		return
	}

	isSuccess, code := auth_service.CheckAuth(params.UserName, params.OldPassword)
	if !isSuccess {
		appG.Response(http.StatusOK, isSuccess, code, nil)
		return
	}

	if ok, err := auth_service.ChangePassword(params.UserName, params.NewPassword); !ok {
		appG.Response(http.StatusUnprocessableEntity, false, e.ERROR_AUTH_CHANGE_PASSWORD_FAIL, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, nil)
	return
}