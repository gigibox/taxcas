package weixin_service

import (
	"taxcas/models"
)

const (
	Col_order  = "order"
	Col_refund = "refund"
)

func Add(order models.WXPayNotifyReq, col string) (bool, error) {
	return models.MgoInsert(order, col)
}
