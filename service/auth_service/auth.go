package auth_service

import "taxcas/models"

func CheckAuth(username, password string) (bool, int) {
	return models.CheckAuth(username, password)
}
