package routers

import (
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"net/http"
	"taxcas/middleware/cors"
	"taxcas/routers/api/admin"
	"taxcas/routers/api/user"
	"taxcas/routers/api/weixin"

	"github.com/gin-gonic/gin"

	_ "taxcas/docs"

	"taxcas/pkg/export"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	// 跨域
	r.Use(cors.Cors())

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	public := r.Group("api")
	{
		public.StaticFS("/export", http.Dir(export.GetExportFullPath()))
		public.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

		public.GET("/admin/login", admin.Login)

		if setting.ServerSetting.RunMode == "debug" {
			public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		public.GET("/e-certs", admin.GetCertsName)
		public.GET("/e-certs/:certid/:id", admin.OfficialWebsite)
	}

	apiv1 := r.Group("/api/v1")
	{
		// 获取字体
		apiv1.GET("/admin/fonts", admin.GetFonts)

		// 获取证书列表
		apiv1.GET("/admin/certs", admin.GetCertsList)

		// 添加证书
		apiv1.POST("/admin/certs", admin.AddCert)

		// 编辑证书
		apiv1.PUT("/admin/certs", admin.EditCert)

		// 上传证书模板
		apiv1.POST("/admin/images", admin.UploadImage)

		// 上传CSV文件
		apiv1.POST("/admin/excels", admin.UploadExcel)

		// 预览证书模板
		apiv1.POST("/admin/images/certs", admin.PreviewImage)

		// 查询申请信息
		apiv1.GET("/admin/applicants/certs/:certid", admin.GetApplicantList)

		// 执行审核结果
		apiv1.PUT("/admin/applicants/certs/:certid", admin.UpdateApplicants)

		// 导出申请信息
		apiv1.GET("/admin/files/applicants/certs/:certid", admin.ExportApplicants)

		// 修改密码
		apiv1.PUT("/admin/password", admin.ChangePassword)

		// 查看证书
		apiv1.GET("/admin/images/certs/:certid/:openid", admin.UserCertificates)
	}

	// 用户端接口
	{
		// 申请证书
		apiv1.POST("/weixin/applicants/users", user.ApplyForCert)

		// 查询申领信息
		apiv1.GET("/weixin/applicants/users/:certid/:openid", user.GetApplicant)

		// 查询用户信息
		apiv1.GET("/weixin/users/:openid", user.GetUserInfo)

		// 查询证书列表
		apiv1.GET("/weixin/certs", admin.GetCertsList)

		// 获取支付订单
		apiv1.GET("/weixin/wxorder/:openid/:certid", weixin.WXPayUnifyOrderReq)

		// 获取openid
		public.GET("/weixin/openid/:code", weixin.WXGetOpenID)

		// 微信服务端回调
		public.POST("weixin/wxnotify", weixin.WXPayCallback)

		// 申请退款
		//apiv1.GET("weixin/wxrefund/:out_trade_no", weixin.WXPayRefund)

		// 查询退款
		apiv1.GET("weixin/wxquery/:openid/:certid", weixin.WXPayRefundQuery)

		// 发送模板消息
		//apiv1.GET("weixin/wxtemplate/:openid", weixin.WXSendTemplateMsg)

		// 发送文本消息
		//apiv1.GET("weixin/wxtext/:openid", weixin.WXSendText)

		// 查看证书
		apiv1.GET("/weixin/e-certs/:certid/:openid", admin.UserCertificates)
	}

	return r
}
