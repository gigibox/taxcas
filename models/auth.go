package models

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"taxcas/pkg/setting"
	"taxcas/pkg/util"
)

func CheckAuth(username, password string) (bool, error) {
	result := C_admin{}

	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + "admin")

	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	if result.Password == util.EncodeMD5(password) {
		return true, nil
	}

	return false, nil
}
