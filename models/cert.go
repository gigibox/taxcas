package models

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"taxcas/pkg/setting"
)

func GetAllCert(collection string) ([]C_certs) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + collection)

	query := c.Find(nil)
	certs := []C_certs{}
	query.All(&certs)

	return certs
}

func CheckCertExist(certname, collection string) (bool, error) {
	result := C_certs{}

	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + collection)

	err := c.Find(bson.M{"certname": certname}).One(&result)
	if err != nil {
		return true, nil
	}

	return false, nil
}

func CountCollection(collection string) (int, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + collection)

	i, err := c.Count()

	return i, err
}

func AddCert(cert C_certs, collection string) (bool, error) {
	c := session.DB(setting.DatabaseSetting.Name).C(setting.DatabaseSetting.TablePrefix + collection)

	err := c.Insert(cert)
	if err != nil {
		log.Println(err)
		return false, nil
	}

	return true, nil
}