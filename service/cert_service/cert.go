package cert_service

import "C"
import (
	"log"
	"strconv"
	"taxcas/models"
)

type S_cert struct {
	Collection string
	Data models.C_certs
}

const (
	col_certs = "certs"
)

func (this *S_cert) CheckExist() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certname", this.Data.CertName, this.Collection, &result)
}

func (this *S_cert) Add() (bool, error) {
	var err error

	this.Data.Status = "enable"
	this.Data.CertID, err = models.MgoCountCollection(this.Collection)

	// 证书模板设计, 暂时写死
	this.Data.ImageDesign.Name = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.EnglishName = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.PersonalID = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.SerialNumber = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.Date = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}

	// 序号从1开始
	this.Data.CertID++
	if err != nil {
		log.Println(err)
		return false, err
	}

	return models.MgoInsert(this.Data, this.Collection)
}

func GetAllCertName() (interface{}) {
	type simpleCert struct {
		Id string `json:"cert_id"`
		Name string `json:"cert_name"`
		Status string `json:"status"`
	}

	certs := []simpleCert{}
	result := []models.C_certs{}

	models.MgoFindAll(col_certs, &result)
	for i, _ := range result {
		t := simpleCert{Id:strconv.Itoa(result[i].CertID), Name:result[i].CertName, Status:result[i].Status}
		certs = append(certs, t)
	}

	return certs
}

func CheckExistByID(id int) (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certid", id, col_certs, &result)
}

func SignCertImage(design models.ImageDesigner) (string, error) {
	imageSaveUrl, err := models.SignImage(design)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return imageSaveUrl, nil
}
