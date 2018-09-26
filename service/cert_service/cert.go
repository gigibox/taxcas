package cert_service

import "C"
import (
	"log"
	"strconv"
	"taxcas/models"
	"taxcas/pkg/export"
	"taxcas/pkg/logging"
	"taxcas/pkg/upload"
	"taxcas/pkg/util"
	"taxcas/service/apply_service"
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
	this.Data.Status = "enable"
	count, err := models.MgoCountCollection(this.Collection)
	if err != nil {
		log.Println(err)
		return false, err
	}
	count++

	// 序号从1开始
	this.Data.CertID = strconv.Itoa(count)

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
		t := simpleCert{Id:result[i].CertID, Name:result[i].CertName, Status:result[i].Status}
		certs = append(certs, t)
	}

	return certs
}

func CheckExistByID(id string) (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certid", id, col_certs, &result)
}

func GetCertByID(id string, doc *models.C_certs) (bool, error) {
	return models.MgoCheckKeyExist("certid", id, col_certs, &doc)
}

func GetCertImage(design *models.ImageDesigner, certid, wechatid string) (string, error) {
	var image string

	// 预览证书, 设计模板
	if design != nil {
		image = upload.GetImagePath() + util.GetRandomFileName("image.png")

		err := models.SignImage(image, design)
		if err != nil {
			log.Println(err)
			return "", err
		}
	} else { // 获取/生成用户证书
		apply := models.C_Apply{}
		cert  := models.C_certs{}
		GetCertByID(certid, &cert)
		apply_service.GetApplyByOpenid(certid, wechatid, &apply)

		// 已存在
		if apply.ImageSaveUrl != "" {
			return apply.ImageSaveUrl, nil
		}

		// 生成
		if apply.CertID == cert.CertID && apply.ApplyStatus == models.Passed && cert.ImageDesign.ImgName != ""{
			designer := cert.ImageDesign
			designer.Name.Str			= apply.Name
			designer.EnglishName.Str 	= apply.EnglishName
			designer.PersonalID.Str		= apply.PersonalID
			designer.SerialNumber.Str	= apply.SerialNumber
			designer.Date.Str			= apply.StudyDate

			image = export.GetExportImagePath(certid) + apply.SerialNumber + ".png"

			// 检查并创建文件夹
			util.CheckDir(export.GetRuntimePath() + export.GetExportImagePath(certid))

			if err := models.SignImage(image, &designer); err != nil {
				log.Println(err)
				return "", err
			}

			// 更新电子证书图片信息
			apply.ImageSaveUrl = image
			if ok , err := models.MgoUpsert("applicant.user.wechatid", wechatid, "cert" + certid + "_apply", apply); !ok {
				logging.Warn("Update applicant status:", err)
			}
		}
	}

	return image, nil
}
