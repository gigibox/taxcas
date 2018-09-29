package weixin

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/objcoding/wxpay"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"sort"
	"strings"
	"taxcas/models"
	"taxcas/pkg/app"
	"taxcas/pkg/e"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"
	"taxcas/pkg/util"
	"taxcas/service/apply_service"
	"taxcas/service/weixin_service"
)

// @Summary 获取用户openid
// @Tags 	微信公众号
// @Param   code path string true "微信浏览器获取的 code"
// @Success 200 {string} json "{"openid":string, "token":string}"
// @Router  /api/v1/weixin/openid/{code} [get]
func WXGetOpenID(c *gin.Context) {
	appG := app.Gin{c}
	type Response struct {
		Access_token  string `json:"access_token"`
		Expires_in    int    `json:"expires_in"`
		Refresh_token string `json:"refresh_token"`
		Openid        string `json:"openid"`
		Scope         string `json:"scope"`
	}
	code := c.Param("code")

	url := strings.Join([]string{"https://api.weixin.qq.com/sns/oauth2/access_token",
		"?appid=", setting.WeixinSetting.AppID,
		"&secret=", setting.WeixinSetting.AppSecret,
		"&code=", code,
		"&grant_type=authorization_code"}, "")

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		appG.Response(http.StatusOK, false, e.ERROR, nil)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR, nil)
	}

	fmt.Println("+++++++++++++++++++++++++")
	fmt.Println(string(body))
	fmt.Println("+++++++++++++++++++++++++")

	result := Response{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR, nil)
	}

	openid := result.Openid
	fmt.Println("=======================")
	fmt.Println(result)
	fmt.Println("=======================")

	// 返回token
	c.Header("Authorization", util.GenerateToken("weixin", openid))
	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"openid": openid,
	})
}

// @Summary 获取支付订单
// @Tags 	微信公众号
// @Security ApiKeyAuth
// @Param   openid path string true "用户 openid"
// @Param   certid path string true "证书 id"
// @Success 200 {string} json "{"prepay_id":string}"
// @Router  /api/v1/weixin/wxorder/{openid}/{certid} [get]
func WXPayUnifyOrderReq(c *gin.Context) {
	appG := app.Gin{c}
	ip := c.ClientIP()
	openid := c.Param("openid")
	certid := c.Param("certid")

	result := models.C_certs{}
	isExist, err := models.MgoFindOne("certid", certid, "certs", &result)
	if err != nil {
		appG.Response(http.StatusOK, false, e.ERROR_EXIST_CERT_FAIL, nil)
	}
	if isExist == false {
		appG.Response(http.StatusOK, false, e.ERROR_NOT_EXIST_CERT, nil)
	}

	price := result.Price
	fmt.Println(price)

	// 在申请订单中填写支付订单号
	applyService := apply_service.S_Apply{
		Collection: "cert" + certid + "_apply",
	}

	apply_service.GetApplyByOpenid(certid, openid, &applyService.Data)
	// 收费证书
	if price > 0 {
		out_trade_no := UniqueId()
		fmt.Println(ip)
		fmt.Println(result)

		client := wxpay.NewClient(wxpay.NewAccount(setting.WeixinSetting.AppID, setting.WeixinSetting.MchID, setting.WeixinSetting.ApiKey, false))
		params := make(wxpay.Params)
		params.SetString("body", result.CertName).
			SetString("out_trade_no", out_trade_no).
			SetInt64("total_fee", int64(price*100)).
			SetString("spbill_create_ip", ip).
			SetString("notify_url", setting.WeixinSetting.Notify_url).
			SetString("openid", openid).
			SetString("attach", certid).
			SetString("trade_type", "JSAPI")

		p, err := client.UnifiedOrder(params)
		if err != nil {
			appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
		}

		prepay_id := p.GetString("prepay_id")
		appid := p.GetString("appid")
		appG.Response(http.StatusOK, true, e.SUCCESS, map[string]interface{}{
			"prepay_id": prepay_id,
			"appid":     appid,
			"price":     price * 100,
			"apikey":    setting.WeixinSetting.ApiKey,
			"orderid":   out_trade_no,
			"name":      result.CertName,
		})

		applyService.Data.ApplyStatus = models.NotPaid
		applyService.Data.ApplyStatusMsg = models.StatusMsg[models.NotPaid]
		applyService.Data.PayOrder = out_trade_no
		applyService.Data.PayStatus = models.NotPaid
	} else {
		applyService.Data.ApplyStatus = models.Pending
		applyService.Data.ApplyStatusMsg = models.StatusMsg[models.Pending]

		appG.Response(http.StatusOK, true, e.SUCCESS, map[string]int{
			"price": price,
		})
	}

	// 更新状态
	applyService.Data.WechatID = openid
	applyService.UpdateStatus()
	fmt.Println("lsdkjgdslj")
	return
}

func WXPayCallback(c *gin.Context) {
	var mr models.WXPayNotifyReq
	var ms models.WXPayNotifyResp
	err := c.Bind(&mr)
	if err != nil {
		ms.Return_code = "FAIL"
		ms.Return_msg = "failed to receive args"
		c.XML(http.StatusOK, ms)
		return
	}

	var reqMap map[string]interface{}
	reqMap = make(map[string]interface{}, 0)

	reqMap["return_code"] = mr.Return_code
	reqMap["return_msg"] = mr.Return_msg
	reqMap["appid"] = mr.Appid
	reqMap["mch_id"] = mr.Mch_id
	reqMap["nonce_str"] = mr.Nonce
	reqMap["result_code"] = mr.Result_code
	reqMap["openid"] = mr.Openid
	reqMap["is_subscribe"] = mr.Is_subscribe
	reqMap["trade_type"] = mr.Trade_type
	reqMap["bank_type"] = mr.Bank_type
	reqMap["total_fee"] = mr.Total_fee
	reqMap["fee_type"] = mr.Fee_type
	reqMap["cash_fee"] = mr.Cash_fee
	reqMap["cash_fee_type"] = mr.Cash_fee_type
	reqMap["transaction_id"] = mr.Transaction_id
	reqMap["out_trade_no"] = mr.Out_trade_no
	reqMap["attach"] = mr.Attach
	reqMap["time_end"] = mr.Time_end

	if wxpayVerifySign(reqMap, mr.Sign) {
		//这里就可以更新我们的后台数据库了，其他业务逻辑同理。
		result, err := weixin_service.Add(mr, weixin_service.Col_order)
		if err != nil {
			fmt.Println("写数据库失败")
		}
		if result {
			applyService := apply_service.S_Apply{
				Collection: "cert" + mr.Attach + "_apply",
			}

			apply_service.GetApplyByOpenid(mr.Attach, mr.Openid, &applyService.Data)
			applyService.Data.PayAmount = mr.Total_fee / 100
			applyService.Data.PayTime = time.Now().Unix()
			applyService.Data.ApplyStatus = models.Pending
			applyService.Data.PayStatus = models.Paid
			applyService.Data.ApplyStatusMsg = models.StatusMsg[models.Paid]
			applyService.UpdateStatus()

			ms.Return_code = "SUCCESS"
			ms.Return_msg = "OK"
		}
	} else {
		ms.Return_code = "FAIL"
		ms.Return_msg = "failed to verify sign, please retry!"
	}

	c.XML(http.StatusOK, ms)
}

func WXPayRefund(out_trade_no string) (bool, error) {
	out_refund_no := UniqueId()
	//通过订单号去支付成功的数据库表中查找是否有此订单，并取出相应的total_fee,设置refund_fee

	account := wxpay.NewAccount(setting.WeixinSetting.AppID, setting.WeixinSetting.MchID, setting.WeixinSetting.ApiKey, false)
	account.SetCertData(upload.GetApiCertFullPath())

	result := models.WXPayNotifyReq{}
	isExist, err := models.MgoFindOne("out_trade_no", out_trade_no, weixin_service.Col_order, &result)
	if err != nil {
		return false, err
	}
	if isExist == false {
		return false, errors.New("此订单不存在")
	}

	client := wxpay.NewClient(account)
	params := make(wxpay.Params)
	params.SetString("out_trade_no", out_trade_no).
		SetString("out_refund_no", out_refund_no).
		SetInt64("total_fee", int64(result.Total_fee)).
		SetInt64("refund_fee", int64(result.Total_fee))

	_, err = client.Refund(params)
	if err != nil {
		return false, err
	}
	return true, nil
}

// @Summary 查询退款
// @Tags 	微信公众号
// @Security ApiKeyAuth
// @Param   openid path string true "用户openid"
// @Param   certid path string true "证书id"
// @Success 200 {string} json "{"msg":string, "extra":}"
// @Router  /api/v1/weixin/wxquery/{openid}/{certid} [get]
func WXPayRefundQuery(c *gin.Context) {
	appG := app.Gin{c}
	openid := c.Param("openid")
	certid := c.Param("certid")

	result := models.C_Apply{}
	apply_service.GetApplyByOpenid(certid, openid, &result)
	if result.PayOrder == "" {
		appG.Response(http.StatusOK, false, e.ERROR, nil)
		return
	}

	out_trade_no := result.PayOrder

	client := wxpay.NewClient(wxpay.NewAccount(setting.WeixinSetting.AppID, setting.WeixinSetting.MchID, setting.WeixinSetting.ApiKey, false))
	params := make(wxpay.Params)
	params.SetString("out_trade_no", out_trade_no)

	p, err := client.RefundQuery(params)
	if err != nil {
		appG.Response(http.StatusOK, false, e.SUCCESS, err)
	}

	appG.Response(http.StatusOK, true, e.SUCCESS, p)
}

func wxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	signCalc := wxpayCalcSign(needVerifyM, setting.WeixinSetting.ApiKey)

	//	slog.Debug("计算出来的sign: %v", signCalc)
	//	slog.Debug("微信异步通知sign: %v", sign)
	if sign == signCalc {
		fmt.Println("签名校验通过!")
		return true
	}

	fmt.Println("签名校验失败!")
	return false
}

func wxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	fmt.Println("微信支付签名计算, API KEY:", key)
	//STEP 1, 对key进行升序排序.
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}

	sort.Strings(sorted_keys)

	//STEP2, 对key=value的键值对用&连接起来，略过空值
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//STEP3, 在键值对的最后加上key=API_KEY
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//STEP4, 进行MD5签名并且将所有字符转为大写.
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func UniqueId() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

/*type AccessToken struct {
	Token      string `json:"access_token"`
	ExpiresIn  int    `json:"expires_in"`
	CreateTime int64
}

type sendTemplateData struct {
	Touser     string                       `json:"touser"`
	TemplateId string                       `json:"template_id"`
	Url        string                       `json:"url"`
	Data       map[string]map[string]string `json:"data"`
}

var accessToken *AccessToken

func WXSendTemplateMsg(c *gin.Context) {
	appG := app.Gin{c}
	openid := c.Param("openid")
	s, _ := sendTmplete(openid, "PhRgTu9MHjAJBSKvNJn9O2t2BNgoAG1z0GvIfUlFQVo", "https://tax.caishuidai.com/")
	b, err := json.Marshal(s)
	if err != nil {
		// fmt.Fprintf(, "error")
		fmt.Println("error")
		appG.Response(http.StatusOK, false, e.ERROR, "发送模板消息失败")
	}
	//返回发送结果
	// fmt.Fprintf(, string(b))
	fmt.Println(string(b))
	appG.Response(http.StatusOK, true, e.SUCCESS, b)
}

func sendTmplete(openid, templateID, backUrl string) (string, error) {
	var templateResult map[string]map[string]map[string]string
	//解析模板消息
	sendTemplateData := new(sendTemplateData)
	sendTemplateData.Touser = openid
	sendTemplateData.TemplateId = templateID
	sendTemplateData.Url = backUrl

	templateData, err := ioutil.ReadFile(upload.GetTemplateFullPath())
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	json.Unmarshal([]byte(templateData), &templateResult)

	sendTemplateData.Data = templateResult[templateID]
	reqBody, err := json.Marshal(sendTemplateData)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	token := getAccessToken()
	url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + token
	res, err := postReq(url, string(reqBody))
	return res, err
}*/

/**
* @note 发起post请求
 * @params url string 请求url
  * @params reqBody string 请求体
   * @return 请求响应
*/
/*func postReq(url, reqBody string) (string, error) {
	//创建请求
	req, err := http.NewRequest("POST", url, strings.NewReader(reqBody))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	//增加header
	req.Header.Set("Content-Type", "application/json; encoding=utf-8")

	//执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("POST请求:创建请求失败", err)
		return "", err
	}
	//读取响应
	body, err := ioutil.ReadAll(resp.Body) //此处可增加输入过滤
	if err != nil {
		fmt.Println("POST请求:读取body失败", err)
		return "", err
	}

	fmt.Println("POST请求:创建成功", string(body))

	defer resp.Body.Close()

	return string(body), nil
}*/

/**
* @note 获取微信的access_token
 * @params *accessToken 用来存放token的结构体
  * @return token string
*/
/*func getAccessToken() string {
	url := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + setting.WeixinSetting.AppID + "&secret=" + setting.WeixinSetting.AppSecret
	if accessToken == nil || (time.Now().Unix()-accessToken.CreateTime) > 7000 {
		accessToken = new(AccessToken)
		res, err := postReq(url, "")

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		var t map[string]interface{}
		json.Unmarshal([]byte(res), &t)
		accessToken.Token = fmt.Sprintf("%s", t["access_token"])
		accessToken.ExpiresIn, _ = fmt.Printf("%d", t["expires_in"])
		accessToken.CreateTime = time.Now().Unix()
	}

	fmt.Println("当前Token为：", accessToken.Token)
	fmt.Printf("过期时间还有：%d秒\r\n", 7200-(time.Now().Unix()-accessToken.CreateTime))
	return accessToken.Token
}

type CustomServiceMsg struct {
	ToUser  string         `json:"touser"`
	MsgType string         `json:"msgtype"`
	Text    TextMsgContent `json:"text"`
}

type TextMsgContent struct {
	Content string `json:"content"`
}

func WXSendText(openid, msg string) (bool, error) {
	accessToken := getAccessToken()
	err := pushCustomMsg(accessToken, openid, msg)
	if err != nil {
		fmt.Println("Push custom service message err:", err)
		return false, err
	}
	return true, nil
}

func pushCustomMsg(accessToken, toUser, msg string) error {
	csMsg := &CustomServiceMsg{
		ToUser:  toUser,
		MsgType: "text",
		Text:    TextMsgContent{Content: msg},
	}

	body, err := json.MarshalIndent(csMsg, " ", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(body))

	postReq, err := http.NewRequest("POST",
		strings.Join([]string{"https://api.weixin.qq.com/cgi-bin/message/custom/send", "?access_token=", accessToken}, ""),
		bytes.NewReader(body))
	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/json; encoding=utf-8")

	client := &http.Client{}
	resp, err := client.Do(postReq)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}*/
