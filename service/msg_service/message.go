package msg_service

import (
	"bytes"
	"fmt"
	"taxcas/models"
	"taxcas/pkg/logging"
	"taxcas/service/weixin_service"
	"text/template"
	"time"
)

const (
	RefundMsg = "{{.Name}}你好，很抱歉你于{{.Date}}提交的{{.Cert}}电子证书领取申请未审核通过，您支付的{{.Amount}}元已为您发起退款，可在【电子证书】中选择证书后查看退款进度并重新申请证书"
	PassedMsg = "{{.Name}}你好，你于{{.Date}}提交的{{.Cert}}电子证书领取申请已审核通过，可在【电子证书】中选择证书查看下载"
	RejectMsg = "{{.Name}}你好，很抱歉你于{{.Date}}提交的{{.Cert}}电子证书领取申请未审核通过，可在【电子证书】中重新申请"
)

var Msg = map[int]string{
	models.Refunded: RefundMsg,
	models.Passed: PassedMsg,
	models.Reject: RejectMsg,
}

type WxMessage struct {
	Name string
	Date string
	Cert string
	Amount int
	Openid string
}

func Send(data *models.C_Apply) {
	if Msg[data.ApplyStatus] == "" || (data.ApplyStatus == models.Reject && data.PayStatus != models.Refunded) { return }

	t := time.Unix(data.ApplyDate, 0)
	m := WxMessage{
		Name: data.Name,
		Cert: data.CertName,
		Date: fmt.Sprintf("%d年%d月%d日", t.Year(), t.Month(), t.Day()),
		Amount: data.PayAmount,
		Openid: data.WechatID,
	}

	var buff bytes.Buffer

	tmpl, err := template.New("weixin").Parse(Msg[data.ApplyStatus])
	if err != nil {
		logging.Warn(err)
	}
	err = tmpl.Execute(&buff, m)
	if err != nil {
		logging.Warn(err)
	}

	weixin_service.WXSendText(data.WechatID, buff.String())
}