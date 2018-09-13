package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"taxcas/pkg/setting"
)

var session *mgo.Session

func init() {
	var err error

	session, err = mgo.Dial(fmt.Sprintf("mongodb://%s", setting.DatabaseSetting.Host))
	//defer session.Close()
	if err != nil {
		log.Println(err)
	}

	session.SetMode(mgo.Monotonic, true)
}