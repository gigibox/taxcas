package routers

import (
	"fmt"
	"net/http"
	"strings"
	"taxcas/routers/api/admin"
	"taxcas/routers/api/user"

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
	// 跨域
	r.Use(Cors())

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

	r.StaticFS("/export", http.Dir(export.GetExportFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//r.StaticFS("/qrcode", http.Dir(qrcode.GetQrCodeFullPath()))

	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		apiv1.GET("/admin/fonts", admin.GetFonts)
		apiv1.GET("/admin/certs", admin.GetCertList)
		apiv1.POST("/admin/cert/add", admin.AddCert)
		apiv1.POST("/admin/cert/upload_image", admin.UploadImage)
		apiv1.GET("/admin/cert/preview", admin.PreviewImage)
		apiv1.GET("/admin/cert/applicants/:cert_id/:action", admin.GetApplicantList)
		apiv1.POST("/admin/export/cert/:cert_id/:action", admin.ExportApplicants)
		apiv1.POST("/admin/import/cert/:cert_id/:action", admin.ImportApplicants)
		/*
		// 导入证书审核结果
		apiv1.POST("/admin/import",)
		// 更新申请状态
		apiv1.PUT("")
		// 删除指定证书
		apiv1.DELETE("/admin")
		*/
	}

	// 用户端接口
	{
		// 提交申请
		apiv1.PUT("/cert/:cert_id/apply", apply.Apply)
	}

	return r
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			//c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}

		c.Next()
	}
}