package admin

import (
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"taxcas/models"
	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/export"
	"taxcas/pkg/logging"
	"taxcas/pkg/upload"
	"taxcas/service/apply_service"
	"taxcas/service/cert_service"
)

// @Summary 	获取证书列表
// @Tags 		后台管理
// @Description 查询所有证书列表
// @Produce  	json
// @Success 	200 {object} app.ResponseMsg "data:[{"cert_id":"0", "cert_name":"证书1", "status":"enable"}]"
// @Router 		/api/v1/admin/certs [get]
func GetCertList(c *gin.Context) {
	appG := app.Gin{c}
	appG.Response(http.StatusOK, true, e.SUCCESS, cert_service.GetAllCertName())
}

// @Summary 	查询证书申领信息
// @Tags 		后台管理
// @Description 查询指定证书的申领信息,
// @Param   	certid path string true "Cert ID"
// @Param   	type query int true "类型 all | export | verify | passed | Reject"
// @Param   	page query int false "页数"
// @Param   	limit query int false "每页显示的条数"
// @Produce  	json
// @Success 	200 {object} app.ResponseMsg "data:[{""}]"
// @Router 		/api/v1/admin/applicants/certs/{certid} [get]
func GetApplicantList(c *gin.Context) {
	appG := app.Gin{c}

	id		:= com.StrTo(c.Query("certid")).MustInt()
	page 	:= com.StrTo(c.Query("page")).MustInt()
	limit	:= com.StrTo(c.Query("limit")).MustInt()
	action	:= c.Param("type")

	if notExist, _ := cert_service.CheckExistByID(id); notExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, nil)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, apply_service.GetApplyList(c.Param("certid"), action, page, limit))
}

// @Summary 添加证书
// @Tags 	后台管理
// @Produce json
// @Param   certInfo body models.C_certs true "证书详细信息"
// @Success 200 {object} app.ResponseMsg "certID 不需要填写, 失败返回 false 及 msg"
// @Router  /api/v1/admin/certs [post]
func AddCert(c *gin.Context) {
	appG := app.Gin{c}

	var cert models.C_certs
	var err error

	contentType := c.Request.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		err = c.BindJSON(&cert)
	case "application/x-www-form-urlencoded":
		err = c.MustBindWith(&cert, binding.FormPost)
	}
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, err)
		return
	}

	valid := validation.Validation{}

	if ok, _ := valid.Valid(&cert); !ok {
		app.MarkErrors(valid.Errors)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, "valid Errors")
		return
	}

	certService := cert_service.S_cert{Collection: "certs", Data: cert}

	if isExist, _ := certService.CheckExist(); !isExist {
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT, nil)
		return
	}

	isAdded, err := certService.Add()

	appG.Response(http.StatusOK, isAdded, e.SUCCESS, err)
}

// @Summary 预览证书
// @Tags 	后台管理
// @Produce json
// @Param   positions body models.ImageDesigner false "证书详细信息"
// @Success 200 {object} app.ResponseMsg "data:{"image_save_path":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router  /api/v1/admin/images/certs [get]
func PreviewImage(c *gin.Context) {
	appG := app.Gin{c}

	var design models.ImageDesigner

	if err := c.Bind(&design); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, err)
		return
	}

	design.Name.Str = "李雷"
	design.EnglishName.Str = "LiLei"
	design.PersonalID.Str = "110010201010010101"
	design.SerialNumber.Str = "2018091401012345"
	design.Date.Str = "2018 年 9 月 14 日"
	// 暂时写死
	design.Name = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	design.EnglishName = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	design.PersonalID = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	design.SerialNumber = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	design.Date = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	imageName, err := cert_service.SignCertImage(design)
	if err != nil {
		appG.Response(http.StatusOK, true, e.ERROR_UPLOAD_CREATE_IMAGE_FAIL, map[string]string{
			"imageURL":      upload.GetImageFullUrl(imageName),
			"imageSavePath": upload.GetImagePath() + imageName,
		})

		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"imageURL":      upload.GetImageFullUrl(imageName),
		"imageSavePath": upload.GetImagePath() + imageName,
	})
}

// @Summary  查询字体列表
// @Tags 	 后台管理
// @Produce  json
// @Success  200 {object} app.ResponseMsg "data:{"image_save_path":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router   /api/v1/admin/fonts [get]
func GetFonts(c *gin.Context) {
	appG := app.Gin{C: c}

	appG.Response(http.StatusOK, true, e.SUCCESS, models.GetFontsList())
}

// @Summary  上传证书模板
// @Tags 	 后台管理
// @Produce  json
// @Param    image formData file true "证书模板图片"
// @Success  200 {object} app.ResponseMsg "data:{"image_save_path":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router   /api/v1/admin/images [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, err)
		return
	}

	if image == nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetRandomFileName(image.Filename)
	fullPath  := upload.GetImageFullPath()

	if !upload.CheckImageExt(imageName) || !upload.CheckFileSize(file) {
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	if err := upload.CheckDir(fullPath); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, err)
		return
	}

	if err := c.SaveUploadedFile(image, fullPath + imageName); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, err)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"imageURL":      upload.GetImageFullUrl(imageName),
		"imageSavePath": upload.GetImagePath() + imageName,
	})
}

// @Summary  导入文件
// @Tags 	 后台管理
// @Param    excel formData file true "审核结果.csv"
// @Success  200 {object} app.ResponseMsg "data:{""}"
// @Router   /api/v1/admin/excels [post]
func UploadExcel(c *gin.Context) {
	appG := app.Gin{C: c}

	file, excel, err := c.Request.FormFile("excel")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, err)
		return
	}

	if excel == nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	saveName := upload.GetRandomFileName(excel.Filename)
	fullPath := upload.GetExcelFullPath()

	if !upload.CheckExcelExt(saveName) || !upload.CheckFileSize(file) {
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_FILE_FORMAT, nil)
		return
	}

	if err := upload.CheckDir(fullPath); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_FILE_FAIL, err)
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
// @Produce  json
// @Param    certid path string true "Cert ID"
// @Param    type query string true "类型 export | Reject"
// @Success  200 {object} app.ResponseMsg "data:{"file_save_path":"upload/images/96a.csv", "file_url": "http://..."}"
// @Router   /api/v1/admin/files/applicants/certs/{certid} [get]
func ExportApplicants(c *gin.Context) {
	appG := app.Gin{c}

	id := com.StrTo(c.Query("certid")).MustInt()

	if isExist, _ := cert_service.CheckExistByID(id); isExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, nil)
		return
	}

	filename, _ := apply_service.ExportFile(c.Param("certid"), c.Param("type"))
	if filename == "" {
		appG.Response(http.StatusOK, false, e.ERROR_EXPORT_FILE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"file_url":      export.GetExportFullUrl(filename),
		"file_save_path": export.GetExportPath() + filename,
	})
}

// @Summary  执行审核结果
// @Tags 	 后台管理
// @Produce  json
// @Param    file query string true "导入的csv文件路径"
// @Param    type query string true "审核中: passed, 已拒绝: refunded"
// @Success  200 {object} app.ResponseMsg "data:{""}"
// @Router   /api/v1/admin/applicants/certs/{certid} [put]
func UpdateApplicants(c *gin.Context) {
	appG := app.Gin{C: c}

	// 检查证书id是否存在
	id := com.StrTo(c.Query("certid")).MustInt()

	if isExist, _ := cert_service.CheckExistByID(id); isExist {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, nil)
		return
	}

	// 解析审核结果
	apply_service.UpdateApplicantsByFile(c.Param("certid"), c.Param("type"), c.Param("file"))

	appG.Response(http.StatusOK, true, e.SUCCESS, nil)
}