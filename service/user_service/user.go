package user_service

import (
	"log"
	"taxcas/models"
)

const (
	col_users = "users"
)

func Add(user models.User) (bool, error){
	doc := models.C_users{
		User: user,
	}

	return models.MgoInsert(doc, "user")
}

func UpdateCerts(user models.User, cert_id string, apply_status int) (bool, error) {
	doc := models.C_users{}
	log.Println(user, "|", cert_id, "|", apply_status)
	if isExist, _ := models.MgoFindOne("user.wechatid", user.WechatID, col_users, &doc); !isExist {
		doc.User = user
		doc.Certs = append(doc.Certs, map[string]int{cert_id:apply_status})

		return models.MgoInsert(doc, col_users)
	} else {
		isAdded := false
		for i := range doc.Certs {
			for k, _ := range doc.Certs[i] {
				if k == cert_id {
					doc.Certs[i][k] = apply_status
					isAdded = true
				}
			}
		}

		if !isAdded {
			doc.Certs = append(doc.Certs, map[string]int{cert_id:apply_status})
		}

		return models.MgoUpsert("user.wechatid", user.WechatID,  col_users, doc)
	}


}

func GetUser(openid string, doc *models.C_users) (bool, error) {

	return models.MgoFindOne("user.wechatid", openid, col_users, &doc)
}