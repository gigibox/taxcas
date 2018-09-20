package apply_service

import (
	"taxcas/models"
	"time"
)

type S_Apply struct {
	Collection string
	Data models.C_Apply
}

func (this *S_Apply) CheckApplyExist() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("applicant.user.wechatid", this.Data.WechatID, this.Collection, &result)
}

func (this *S_Apply) CheckCertByName() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certname", this.Data.CertName, "certs", &result)
}

func (this *S_Apply) CheckApplyStatus() (bool, error) {
	result := models.C_certs{}
	models.MgoCheckKeyExist("certname", this.Data.CertName, "certs", &result)

	if result.Status == "enable" {
		return true, nil
	}

	return false, nil
}

func (this *S_Apply) Add() (bool, error) {
	return models.MgoInsert(this.Data, this.Collection)
}

func New(col string, commit models.Applicant) (*S_Apply) {
	return &S_Apply{
		Collection: col,
		Data: models.C_Apply{
			Applicant: commit,
			ApplyDate: int(time.Now().Unix()),
		},
	}
}
