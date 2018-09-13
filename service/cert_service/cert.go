package cert_service

import "C"
import (
	"log"
	"taxcas/models"
)

type S_cert struct {
	Collection string
	Data models.C_certs
}

func GetAllCertName() (interface{}) {
	type certName struct {
		Id int
		Name string
	}

	certs := []certName{}

	r := models.GetAllCert("certs")
	for i, _ := range r {
		t := certName{Id:r[i].CertID, Name:r[i].CertName}
		certs = append(certs, t)
	}

	return certs
}

func (c *S_cert) CheckExist() (bool, error) {
	return models.CheckCertExist(c.Data.CertName, c.Collection)
}

func (c *S_cert) Add() (bool, error) {
	var err error

	c.Data.CertID, err = models.CountCollection(c.Collection)

	if err != nil {
		log.Println(err)
		return false, err
	}

	return models.AddCert(c.Data, c.Collection)
}