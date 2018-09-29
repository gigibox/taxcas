package models

const PassedMsg = "xxx你好，很抱歉你于xx年xx月xx日提交的xxx电子证书领取申请未审核通过，您支付的xx元已为您发起退款，可在【电子证书】中选择证书后查看退款进度并重新申请证书"

type Inventory struct {
	Material string
	Count    uint
}