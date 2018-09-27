package util

import (
	"github.com/dgrijalva/jwt-go"
	"taxcas/pkg/logging"
	"taxcas/pkg/setting"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	User string `json:"user"`
	Permission string `json:"permission"`
	RefeshTime int64 `json:"ref"`
	jwt.StandardClaims
}

func getExpireTime(rols string) (int64, int64) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)
	RefeshTime := nowTime.Add(1 * time.Hour)

	switch rols {
		case "admin":
			expireTime = nowTime.Add(10 * time.Minute)
			RefeshTime = nowTime.Add(5 * time.Minute)
	}

	return expireTime.Unix(), RefeshTime.Unix()
}

func GenerateToken(roles, user string) (string) {
	exp, ref := getExpireTime(roles)

	claims := Claims{
		EncodeMD5(user),
		roles,
		ref,
		jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			ExpiresAt: exp,
			Issuer:    "taxcas-caishuidai",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		logging.Error(err)
	}

	return token
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func RefreshToken(token string) (string) {
	claims, _ := ParseToken(token)

	// 从redis读取已刷新的token, 若不存在, 生成新的

	return GenerateToken(claims.User, claims.Permission)
}