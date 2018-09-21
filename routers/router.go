package routers

import (
	"net/http"
	"taxcas/middleware/cors"
	"taxcas/routers/api/admin"
	"taxcas/routers/api/user"
	"taxcas/routers/api/weixin"

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

	// 跨域
	r.Use(cors.Cors())

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	gin.SetMode(setting.ServerSetting.RunMode)

	public := r.Group("api")
	{
		public.StaticFS("/export", http.Dir(export.GetExportFullPath()))
		public.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

		public.GET("/auth", api.GetAuth)
		public.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		// 获取字体
		apiv1.GET("/admin/fonts", admin.GetFonts)

		// 获取证书列表
		apiv1.GET("/admin/certs", admin.GetCertList)

		// 添加证书
		apiv1.POST("/admin/certs", admin.AddCert)

		// 上传图片
		apiv1.POST("/admin/images", admin.UploadImage)

		// 预览证书
		apiv1.GET("/admin/images/certs", admin.PreviewImage)

		// 获取申请状态
		apiv1.GET("/admin/applicants/certs/:certid", admin.GetApplicantList)

		// 导入审核结果
		apiv1.POST("/admin/applicants/certs/:certid", admin.ImportApplicants)

		// 导出申请状态
		apiv1.GET("/admin/files/applicants/certs/:certid", admin.ExportApplicants)
	}

	// 用户端接口
	{
		// 申请证书
		apiv1.POST("/weixin/applicants/users/:openid", user.ApplyForCert)

		// 查询用户
		apiv1.GET("/weixin/users/:openid", user.GetUserInfo)

		// 查询证书
		apiv1.GET("/weixin/certs", admin.GetCertList)

		// 获取openid
		apiv1.GET("/weixin/openid/:code", weixin.WXGetOpenID)

		// 获取支付订单
		apiv1.GET("/weixin/wxorder/:openid/:certid", weixin.WXPayUnifyOrderReq)

		// 微信服务端回调
		apiv1.GET("weixin/wxnotify", weixin.WXPayCallback)

		//申请退款
		apiv1.GET("weixin/wxrefund/:out_trade_no", weixin.WXPayRefund)

		//查询退款
		apiv1.GET("weixin/wxquery/:out_refund_no", weixin.WXPayRefundQuery)
	}

	return r
}
