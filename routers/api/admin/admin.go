package admin

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"taxcas/models"
	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/logging"
	"taxcas/pkg/upload"
	"taxcas/pkg/util"
	"taxcas/service/apply_service"
	"taxcas/service/cert_service"
)

// @Summary 	获取证书列表
// @Tags 		后台管理
// @Description 查询所有证书列表
// @Security	ApiKeyAuth
// @Produce  	json
// @Success 	200 {object} app.ResponseMsg "data:[{"cert_id":"0", "cert_name":"证书1", "status":"enable"}]"
// @Router 		/api/v1/admin/certs [get]
func GetCertsList(c *gin.Context) {
	appG := app.Gin{c}
	appG.Response(http.StatusOK, true, e.SUCCESS, cert_service.GetCertsList())
}

// @Summary 	查询证书申领信息
// @Tags 		后台管理
// @Description 查询指定证书的申领信息
// @Security	ApiKeyAuth
// @Param   	certid path string true "Cert ID"
// @Param   	type query string true "类型 all | export | verify | passed | Reject"
// @Param   	page query int false "页数"
// @Param   	limit query int false "每页显示的条数"
// @Param   	field query string false "指定条件查询, 姓名或身份证号"
// @Produce  	json
// @Success 	200 {object} app.ResponseMsg "data:[{""}]"
// @Router 		/api/v1/admin/applicants/certs/{certid} [get]
func GetApplicantList(c *gin.Context) {
	appG := app.Gin{c}

	id    := c.Param("certid")
	act   := c.Query("type")
	fie   := c.Query("field")
	page  := com.StrTo(c.Query("page")).MustInt()
	limit := com.StrTo(c.Query("limit")).MustInt()

	if notExist, _ := cert_service.CheckExistByID(id); notExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, nil)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, apply_service.GetApplyList(c.Param("certid"), act, page, limit, fie))
}

// @Summary 添加证书
// @Tags 	后台管理
// @Security	ApiKeyAuth
// @Produce json
// @Param   certInfo body models.C_certs true "证书详细信息"
// @Success 200 {object} app.ResponseMsg "certID 不需要填写, 失败返回 false 及 msg"
// @Router  /api/v1/admin/certs [post]
func AddCert(c *gin.Context) {
	appG := app.Gin{c}

	var cert models.C_certs

	err := c.Bind(&cert)
	if err != nil {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, err)
		return
	}

	valid := validation.Validation{}

	if ok, _ := valid.Valid(&cert); !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, "valid Errors")
		return
	}

	certService := cert_service.S_cert{Collection: "certs", Data: cert}

	if isExist, _ := certService.CheckExist(); !isExist {
		appG.Response(http.StatusConflict, false, e.ERROR_EXIST_CERT, nil)
		return
	}


	isAdded, err := certService.Add()

	appG.Response(http.StatusCreated, isAdded, e.SUCCESS, err)
}

// @Summary 编辑证书
// @Tags 	后台管理
// @Security ApiKeyAuth
// @Produce json
// @Param   certInfo body models.C_certs true "证书详细信息"
// @Success 200 {object} app.ResponseMsg "失败返回 false 及 msg"
// @Router  /api/v1/admin/certs [put]
func EditCert(c *gin.Context) {
	appG := app.Gin{c}

	var cert models.C_certs

	err := c.Bind(&cert)
	if err != nil {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, err)
		return
	}

	valid := validation.Validation{}

	if ok, _ := valid.Valid(&cert); !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, "valid Errors")
		return
	}

	certService := cert_service.S_cert{Collection: "certs", Data: cert}
	ok, err := certService.Edit()

	appG.Response(http.StatusOK, ok, e.SUCCESS, err)
}

// @Summary 预览证书
// @Tags 	后台管理
// @Security ApiKeyAuth
// @Produce json
// @Param   positions body models.ImageDesigner false "证书详细信息"
// @Success 200 {object} app.ResponseMsg "data:{"image_save_path":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router  /api/v1/admin/images/certs [post]
func PreviewImage(c *gin.Context) {
	appG := app.Gin{c}

	var cert models.C_certs

	if err := c.Bind(&cert); err != nil {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, err)
		return
	}

	cert.ImageDesign.Name.Str = "李雷"
	cert.ImageDesign.EnglishName.Str = "LiLei"
	cert.ImageDesign.PersonalID.Str = "123456197805043210"
	cert.ImageDesign.SerialNumber.Str = "20181000123456"
	cert.ImageDesign.Date.Str = "2018      10      9"

	image, err := cert_service.GetCertImage(&cert.ImageDesign, nil)
	if err != nil {
		appG.Response(http.StatusUnprocessableEntity, false, e.ERROR_UPLOAD_CREATE_IMAGE_FAIL, nil)

		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"imageURL":      util.GetAppFullUrl(image),
		"imageSavePath": image,
	})
}

// @Summary  查询字体列表
// @Tags 	 后台管理
// @Security ApiKeyAuth
// @Produce  json
// @Success  200 {object} app.ResponseMsg "data:{"image_save_path":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router   /api/v1/admin/fonts [get]
func GetFonts(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, true, e.SUCCESS, util.GetFontsList())
}

// @Summary  上传证书模板
// @Tags 	 后台管理
// @Security ApiKeyAuth
// @Produce  json
// @Param    image formData file true "证书模板图片"
// @Success  200 {object} app.ResponseMsg "data:{"image_save_path":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router   /api/v1/admin/images [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, err)
		return
	}

	if image == nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	imageName := util.GetRandomFileName(image.Filename)
	fullPath  := upload.GetImageFullPath()

	if !upload.CheckImageExt(imageName) || !upload.CheckFileSize(file) {
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	if err := c.SaveUploadedFile(image, fullPath + imageName); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"imageURL":      util.GetAppFullUrl(upload.GetImagePath() + imageName),
		"imageSavePath": upload.GetImagePath() + imageName,
	})
}

// @Summary  导入文件
// @Tags 	 后台管理
// @Security ApiKeyAuth
// @Param    excel formData file true "审核结果.csv"
// @Success  200 {object} app.ResponseMsg "data:{""}"
// @Router   /api/v1/admin/excels [post]
func UploadExcel(c *gin.Context) {
	appG := app.Gin{C: c}

	file, excel, err := c.Request.FormFile("excel")
	if err != nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, err)
		return
	}

	if excel == nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	saveName := util.GetRandomFileName(excel.Filename)
	fullPath := upload.GetExcelFullPath()

	if !upload.CheckExcelExt(saveName) || !upload.CheckFileSize(file) {
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_FILE_FORMAT, nil)
		return
	}

	if err := c.SaveUploadedFile(excel, fullPath + saveName); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_SAVE_FILE_FAIL, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"excel_save_path": fullPath + saveName,
	})
}

// @Summary  导出用户申领信息
// @Tags 	 后台管理
// @Security ApiKeyAuth
// @Produce  json
// @Param    certid path string true "Cert ID"
// @Param    type query string true "类型 export | reject"
// @Success  200 {object} app.ResponseMsg "data:{"file_save_path":"upload/excel/96a.csv", "file_url": "http://..."}"
// @Router   /api/v1/admin/files/applicants/certs/{certid} [get]
func ExportApplicants(c *gin.Context) {
	appG := app.Gin{c}

	certid := c.Param("certid")

	if isExist, err := cert_service.CheckExistByID(certid); isExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, err)
		return
	}

	filename, err := apply_service.ExportFile(certid, c.Query("type"))
	if filename == "" {
		appG.Response(http.StatusOK, false, e.ERROR_EXPORT_FILE_FAIL, err)
		return
	} else if filename == "none" {
		appG.Response(http.StatusOK, false, e.ERROR_EXPORT_EMPYT_FILE, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"file_url":      util.GetAppFullUrl(filename),
		"file_save_path": filename,
	})
}

type parameters struct {
	FilePath    string `json:"file_path"`
	Action      string `json:"action"`
	Pids        []string `json:"pids"`
}

// @Summary  执行审核结果
// @Tags 	 后台管理
// @Security ApiKeyAuth
// @Produce  json
// @Param    certid path string true "证书id"
// @Param    data body admin.parameters true "file_path: 导入的csv文件路径, action: passed(审核中) | refunded(已拒绝) | Reject, pids[]: 选中的记录"
// @Success  200 {object} app.ResponseMsg "data:{""}"
// @Router   /api/v1/admin/applicants/certs/{certid} [put]
func UpdateApplicants(c *gin.Context) {
	appG := app.Gin{C: c}

	params := parameters{}
	err := c.BindJSON(&params)
	if err != nil {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, err)
		return
	}

	// 检查证书id是否存在
	certid := c.Param("certid")

	if isExist, _ := cert_service.CheckExistByID(certid); isExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, nil)
		return
	}

	// 解析审核结果
	s, f := apply_service.UpdateApplicants(certid, params.Action, params.FilePath, params.Pids)

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]int{
		"success" : s,
		"failure" : f,
	})
}

// @Summary 查看用户证书
// @Tags 	后台管理
// @Param   certid path string true "证书id"
// @Param   openid path string true "用户微信id"
// @Success 200 {object} app.ResponseMsg "data:{"image_save_path":"export/images/96a.jpg", "image_url": "http://..."}"
// @Router  /api/v1/admin/images/certs/{certid}/{openid} [get]
func UserCertificates(c *gin.Context) {
	appG := app.Gin{c}

	certid := c.Param("certid")
	openid := c.Param("openid")

	apply := models.C_Apply{}
	if certid == "" || openid== ""{
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, nil)
		return
	}

	apply_service.GetApplyByOpenid(certid, openid, &apply)

	if apply.SerialNumber == "" || apply.ApplyStatus != models.Passed {
		appG.Response(http.StatusOK, false, e.ERROR_GET_USER_CERT_IMAGES, nil)
		return
	}

	image, err := cert_service.GetCertImage(nil, &apply)
	if image == "" {
		appG.Response(http.StatusUnprocessableEntity, false, e.ERROR_GET_USER_CERT_IMAGES, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"image_url":       util.GetAppFullUrl(image),
		"image_save_path": image,
	})
}

// @Summary 查询用户证书
// @Tags 	开放接口
// @Param   certid path string true "证书id"
// @Param   id path string true "证书编号或身份证号"
// @Success 200 {object} app.ResponseMsg "data:{"image_save_path":"export/images/96a.jpg", "image_url": "http://..."}"
// @Router  /api/e-certs/{certid}/{id} [get]
func OfficialWebsite(c *gin.Context) {
	appG := app.Gin{c}

	certid := c.Param("certid")
	id := c.Param("id")

	apply := models.C_Apply{}
	if len(id) == 18 {
		apply_service.GetApplyByPID(certid, id, &apply)
	} else if len(id) == 14 {
		apply_service.GetApplyBySN(certid, id, &apply)
	} else {
		appG.Response(http.StatusBadRequest, false, e.INVALID_PARAMS, nil)
		return
	}

	if apply.SerialNumber == "" || apply.ApplyStatus != models.Passed {
		appG.Response(http.StatusOK, false, e.ERROR_GET_USER_CERT_IMAGES, nil)
		return
	}

	image, err := cert_service.GetCertImage(nil, &apply)
	if image == "" {
		appG.Response(http.StatusUnprocessableEntity, false, e.ERROR_GET_USER_CERT_IMAGES, err)
		return
	}

	pdf, err := cert_service.GetCertFile(&apply)
	if pdf == "" {
		appG.Response(http.StatusUnprocessableEntity, false, e.ERROR_GET_USER_CERT_FILES, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"cert_name":     apply.CertName,
		"cert_sn":       apply.SerialNumber,
		"user_name":     apply.Name,
		"user_id":       apply.PersonalID,
		"pdf_url":       util.GetAppFullUrl(pdf),
		"image_url":     util.GetAppFullUrl(image),
		"pdf_save_path": pdf,
		"image_save_path": image,
	})
}

// @Summary 查询证书列表
// @Tags 	开放接口
// @Success 200 {object} app.ResponseMsg "简单的证书信息, 仅包含ID 名称 "
// @Router  /api/e-certs [get]
func GetCertsName(c *gin.Context) {
	appG := app.Gin{c}
	appG.Response(http.StatusOK, true, e.SUCCESS, cert_service.GetAllCertName())
}


// @Summary 查看用户证书
// @Tags 	微信公众号
// @Param   certid path string true "证书id"
// @Param   openid path string true "用户微信id"
// @Success 200 {object} app.ResponseMsg "data:{"image_save_path":"export/images/96a.jpg", "image_url": "http://..."}"
// @Router  /api/v1/weixin/e-certs/{certid}/{openid} [get]
func swagger_a(c *gin.Context) {
}

// @Summary 	获取证书列表
// @Tags 		微信公众号
// @Security	ApiKeyAuth
// @Description 查询所有证书列表
// @Produce  	json
// @Success 	200 {object} app.ResponseMsg "data:[{"cert_id":"0", "cert_name":"证书1", "status":"enable"}]"
// @Router 		/api/v1/weixin/certs [get]
func swagger_c(c *gin.Context) {

}