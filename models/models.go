package models

import (
	"log"
	"taxcas/pkg/setting"
	"taxcas/pkg/util"
)

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}


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

	administrator := t_admin{
		"admin",
		util.EncodeMD5("admin"),
	}

	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + "admin")
	err = c.Insert(administrator)
	if err != nil {
		log.Println(err)
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
