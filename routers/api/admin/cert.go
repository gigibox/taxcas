package admin

import (
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"taxcas/models"
	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/logging"
	"taxcas/pkg/upload"
	"taxcas/service/cert_service"
)

// @Summary 获取证书列表
// @Produce  json
// @Success 200 {string} json "{"code":200,"success": bool,"data":[]}"
// @Router /api/v1/admin/certs [get]
func GetCertList(c *gin.Context) {
	appG := app.Gin{c}
	appG.Response(http.StatusOK, true, e.SUCCESS, cert_service.GetAllCertName())
}

// @Summary 添加证书
// @Produce  json
// @Success 200 {string} json "{"code":200,"success": bool,"data":[]}"
// @Router /api/v1/admin/cert/add [post]
func AddCert(c *gin.Context) {
	appG := app.Gin{c}

	var cert models.C_certs
	var err error

	err = c.BindJSON(&cert)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	valid := validation.Validation{}

	ok, _ := valid.Valid(&cert)
	if !ok {
		app.MarkErrors(valid.Errors)
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	certService := cert_service.S_cert{Collection: "certs", Data: cert}

	isExist, err := certService.CheckExist()
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT_FAIL, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT, nil)
		return
	}

	isAdded, err := certService.Add()

	appG.Response(http.StatusOK, isAdded, e.SUCCESS, err)
}

// @Summary 预览证书
// @Produce  json
// @Param image post file true "图片文件"
// @Success 200 {string} json "{"code":200,"success": bool,"data":{"image_save_url":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router /api/v1//admin/cert/preview [get]
func PreviewImage(c *gin.Context) {

}

// @Summary 上传图片
// @Produce  json
// @Param image post file true "图片文件"
// @Success 200 {string} json "{"code":200,"success": bool,"data":{"image_save_url":"upload/images/96a.jpg", "image_url": "http://..."}"
// @Router /api/v1/admin/cert/upload_image [post]
func UploadImage(c *gin.Context) {
	appG := app.Gin{C: c}

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR, nil)
		return
	}

	if image == nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err := c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		appG.Response(http.StatusOK, false, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	})
}
