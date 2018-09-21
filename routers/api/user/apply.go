package apply

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"taxcas/models"
	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/logging"
	"taxcas/service/apply_service"
	"taxcas/service/user_service"
)

// @Summary 申请证书
// @Tags 	用户申请
// @Produce json
// @Param   applicant body models.Applicant true "用户提交信息"
// @Success 200 {object} app.ResponseMsg "cost 与 applyStatus 不提交. 失败返回 false 及 msg"
// @Router  /api/v1/weixin/applicants [post]
func Apply(c *gin.Context) {
	appG := app.Gin{c}

	var commit models.Applicant
	var err error

	err = c.BindJSON(&commit)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, "BindJson")
		return
	}

	valid := validation.Validation{}

	ok, _ := valid.Valid(&commit)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, "valid error")
		return
	}

	applyService := apply_service.New("cert" + commit.CertID + "_apply", commit)

	// 判断证书是否存在, 或关闭申请
	isExist, err := applyService.CheckCertByName()
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT_FAIL, err)
		return
	}

	if isExist {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, err)
		return
	}

	isOpen, err := applyService.CheckApplyStatus()
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT_FAIL, err)
		return
	}

	if !isOpen {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_CERT_APPLY_DISABLED, err)
		return
	}

	if isApplied, err := applyService.CheckApplyExist(); !isApplied {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_APPLY, err)
		return
	}

	// 提交申请
	isAdded, err := applyService.Add()
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_ADD_APPLY, err)
		return
	}

	// 更新用户信息
	user_service.UpdateCerts(commit.User, applyService.Data.CertID, applyService.Data.ApplyStatus)
	appG.Response(http.StatusOK, isAdded, e.SUCCESS, err)
}
