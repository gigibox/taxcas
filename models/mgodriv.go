package models

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func MgoInsert(doc interface{}, col string) (bool, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	err := c.Insert(doc)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, nil
}

func MgoCheckKeyExist(key string, value interface{}, col string, result interface{}) (bool, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	err := c.Find(bson.M{key: value}).One(result)
	if err != nil {
		return true, err
	}

	return false, err
}

func MgoUpdate(key, value, col string, doc interface{}) (bool, error) {
	return true, nil
}

func MgoUpsert(key, value, col string, doc interface{}) (bool, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	_, err := c.Upsert(bson.M{key: value}, doc)
	if err != nil {
		return false, err
	}

	return true, err
}

func MgoFindOne(key, value, col string, result interface{}) (bool, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	err := c.Find(bson.M{key: value}).One(result)
	if err != nil {
		return false, err
	}

	return true, err
}

func MgoFind(key string, value interface{}, col string, page, limit int, result interface{}) (bool, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	err := c.Find(bson.M{key: value}).Skip(page * limit).Limit(limit).All(result)
	if err != nil {
		return false, err
	}

	return true, err
}

func MgoFindAll(col string, result interface{}) (error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	return c.Find(nil).All(result)
}

func MgoCountQuery(key string, value interface{}, col string) (int, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col)

	return c.Find(bson.M{key: value}).Count()
}

func MgoCountCollection(col string) (int, error) {
	return session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + col).Count()
}
