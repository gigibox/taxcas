package auth_service

import (
	"taxcas/models"
	"taxcas/pkg/e"
	"taxcas/pkg/util"
)

func CheckAuth(username, password string) (bool, int) {
	doc := models.C_admin{}

	if _, err := models.MgoFindOne(username, password, "admin", &doc); err != nil {
		return false, e.ERROR_AUTH_CHECK_USRNAME_FAIL
	}

	if doc.Password != util.EncodeMD5(password) {
		return false, e.ERROR_AUTH_CHECK_PASSWORD_FAIL
	}

	return true, e.SUCCESS
}

func ChangePassword(username, password string) (bool, error) {
	administrator := models.C_admin{
		username,
		util.EncodeMD5(password),
	}

	if ok, err := models.MgoUpdateAll(username, password, "admin", administrator); !ok {
		return false, err
	}

	return true, nil
}