package user

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
// @Tags 	微信公众号
// @Security ApiKeyAuth
// @Produce json
// @Param   applicant body models.Applicant true "用户提交信息"
// @Success 200 {object} app.ResponseMsg "cost 与 applyStatus 不提交. 失败返回 false 及 msg"
// @Router  /api/v1/weixin/applicants/users [post]
func ApplyForCert(c *gin.Context) {
	appG := app.Gin{c}

	var commit models.Applicant
	var err error
	var ok bool

	err = c.BindJSON(&commit)
	if err != nil {
		logging.Debug(err)
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, "BindJson")
		return
	}

	valid := validation.Validation{}

	ok, _ = valid.Valid(&commit)
	if !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, "valid error")
		return
	}

	applyService := apply_service.New("cert" + commit.CertID + "_apply", commit)

	// 判断证书是否存在, 或关闭申请
	isExist, err := applyService.CheckCertByName()
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT_FAIL, err)
		return
	}

	if isExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, err)
		return
	}

	isOpen, err := applyService.CheckApplyStatus()
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT_FAIL, err)
		return
	}

	if !isOpen {
		appG.Response(http.StatusOK, false, e.ERROR_CERT_APPLY_DISABLED, err)
		return
	}

	// 生成编号
	if !applyService.UpdateSerialNumber() {
		appG.Response(http.StatusOK, false, e.ERROR, "生成证书编号错误")
		return
	}

	// 判断是否是重新申领
	dbdate := models.C_Apply{}
	apply_service.GetApplyByOpenid(commit.CertID, commit.WechatID, &dbdate)

	if dbdate.ApplyStatus != models.Reject {
		// 同一微信号只能申请一次
		if isApplied, err := applyService.CheckApplyExistByWX(); !isApplied {
			appG.Response(http.StatusOK, false, e.ERROR_EXIST_APPLY, err)
			return
		}

		// 同一个身份证只能申请一次
		if isApplied, err := applyService.CheckApplyExistByID(); !isApplied {
			appG.Response(http.StatusOK, false, e.ERROR_EXIST_APPLY_PID, err)
			return
		}

		// 提交申请
		ok, err = applyService.Add()
		if err != nil {
			logging.Warn(err)
			appG.Response(http.StatusOK, false, e.ERROR_ADD_APPLY, err)
			return
		}
	} else {
		if dbdate.PayStatus == models.Paid || dbdate.PayStatus == models.Refunding {
			appG.Response(http.StatusOK, false, e.ERROR_EXIST_APPLY_PAY, err)
			return
		}

		ok, err = applyService.Update()
		if err != nil {
			logging.Warn(err)
			appG.Response(http.StatusOK, false, e.ERROR_ADD_APPLY, err)
			return
		}
	}

	// 更新用户信息
	user_service.UpdateCerts(commit.User, applyService.Data.CertID, applyService.Data.ApplyStatus)
	appG.Response(http.StatusCreated, ok, e.SUCCESS, err)
}

// @Summary 查询用户信息
// @Tags 	微信公众号
// @Security ApiKeyAuth
// @Description "申请状态 0:"未支付", 1:"已支付", 2: "待审核", 3: "审核中", 4:"已拒绝", 5:"已通过", 6:"退款中", 7:"已退款", 8:"错误状态""
// @Param   openid path string true "用户openid"
// @Success 200 {object} app.ResponseMsg "用户基本信息 及 证书申领状态 ["申请证书id" : "申请状态"]"
// @Router  /api/v1/weixin/users/{openid} [get]
func GetUserInfo(c *gin.Context) {
	appG := app.Gin{c}

	user := models.C_users{}

	openid	:= c.Param("openid")

	if ok, err := user_service.GetUser(openid, &user); !ok {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_USER, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, user)
}

// @Summary 查询申领信息
// @Tags 	微信公众号
// @Security ApiKeyAuth
// @Param   certid path string true "证书id"
// @Param   openid path string true "用户openid"
// @Success 200 {object} app.ResponseMsg "用户申请证书所提交的所有信息"
// @Router  /api/v1/weixin/applicants/users/{certid}/{openid} [get]
func GetApplicant(c *gin.Context) {
	appG := app.Gin{c}

	certid := c.Param("certid")
	openid := c.Param("openid")

	apply := models.C_Apply{}

	if apply_service.GetApplyByOpenid(certid, openid, &apply); apply.SerialNumber == "" {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_USER, nil)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, apply)
}