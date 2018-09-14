package routers

import (
	"net/http"
	"taxcas/routers/api/admin"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "taxcas/docs"

	"taxcas/pkg/export"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"
	"taxcas/routers/api"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)

	// html
	r.LoadHTMLGlob("views/*")
	webadmin := r.Group("/admin")
	webadmin.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{"title":"taxcas"})
	})
	webadmin.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title":"taxcas"})
	})

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/admin/certs", admin.GetCertList)
		apiv1.POST("/admin/cert/add", admin.AddCert)
		apiv1.POST("/admin/cert/upload_image", admin.UploadImage)
		apiv1.GET("/admin/cert/preview", admin.PreviewImage)
		/*
		// 获取证书申请信息
		apiv1.GET("/admin/cert/:cert_id", admin.GetCertList)
		// 导出证书申请信息
		apiv1.POST("/admin/export",)
		// 导入证书审核结果
		apiv1.POST("/admin/import",)
		// 更新申请状态
		apiv1.PUT("")
		// 添加证书
		apiv1.POST("")
		// 删除指定证书
		apiv1.DELETE("/admin")
		*/
	}

	return r
}
