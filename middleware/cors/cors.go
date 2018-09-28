package cors

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}

		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = "Authorization, Content-Type, " + headerStr
		} else {
			headerStr = "Authorization, Connection, Content-Length, Accept, User-Agent, Content-Type, Accept-Encoding, Origin, Referer, Accept-Language"
		}

		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", headerStr)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Header("Access-Control-Expose-Headers", "Authorization, Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"msg":  "Options Request!",
			})
		}

		c.Next()
	}
}
