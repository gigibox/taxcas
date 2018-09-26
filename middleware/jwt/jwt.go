package jwt

import (
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

		authString := c.GetHeader("Authorization")

		kv := strings.Split(authString, " ")
		if len(kv) != 2 || kv[0] != "Bearer" {
			code = e.ERROR_AUTH
		} else {
			tokenString := kv[1]
			claims, err := util.ParseToken(tokenString)
			if err != nil {
				code = e.ERROR_AUTH
			} else if !strings.Contains(c.Request.RequestURI, claims.Permission){
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if ((claims.ExpiresAt - claims.IssuedAt) % (claims.ExpiresAt - time.Now().Unix())) > 4 {
				// 四分之一时间后, 刷新token
				c.Header("Authorization", util.RefreshToken(tokenString))
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"msg":  e.GetMsg(code),
				"data": nil,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
