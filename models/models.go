package models

import (
	"io/ioutil"
	"log"
	"taxcas/pkg/setting"
	"taxcas/pkg/util"
)

func Setup() {
	var err error

	// 新建数据库, 初始化管理员
	names, err := session.DatabaseNames()
	if err != nil {
		log.Println(err)
	}

	for i := range names {
		if names[i] == setting.DatabaseSetting.Name {
			return
		}
	}

	administrator := C_admin{
		"admin",
		util.EncodeMD5("admin"),
	}

	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + "admin")
	err = c.Insert(administrator)
	if err != nil {
		log.Println(err)
	}
}

func GetFontsList() []string {
	var fonts []string

	list, err := ioutil.ReadDir(setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath)
	if err != nil {
		log.Println(err)
		return fonts
	}

	for _, v := range list {
		fonts = append(fonts, v.Name())
	}

	return fonts
}
