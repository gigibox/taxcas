package jwt

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"taxcas/pkg/e"
	"taxcas/pkg/util"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		httpCode := http.StatusUnauthorized
		nowTime := time.Now().Unix()

		authString := c.GetHeader("Authorization")
		if authString == "" {
			authString = c.PostForm("Authorization")
		}
		kv := strings.Split(authString, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			code = e.ERROR_AUTH
		} else {
			token := kv[1]
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_TOKEN
				log.Println(err)
			} else if !strings.Contains(c.Request.RequestURI, claims.Permission) {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				httpCode = http.StatusForbidden
			} else if nowTime > claims.RefeshTime {
				c.Header("Authorization", util.RefreshToken(token))
			} else {

			}
		}

		if code != e.SUCCESS {
			c.JSON(httpCode, gin.H{
				"success": false,
				"msg":     e.GetMsg(code),
				"data":    nil,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
