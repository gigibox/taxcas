package models

import (
	"gopkg.in/mgo.v2/bson"
	"taxcas/pkg/e"
	"taxcas/pkg/setting"
	"taxcas/pkg/util"
)

func CheckAuth(username, password string) (bool, int) {
	result := C_admin{}

	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + "admin")

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return false, e.ERROR_AUTH_CHECK_USRNAME_FAIL
	}

	if result.Password != util.EncodeMD5(password) {
		return false, e.ERROR_AUTH_CHECK_PASSWORD_FAIL
	}

	return true, 0
}
