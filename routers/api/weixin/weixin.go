package weixin

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/objcoding/wxpay"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
//	"taxcas/models"
	"taxcas/pkg/app"
	"taxcas/pkg/config"
	"taxcas/pkg/e"
	//	"errors"
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
		"?appid=", config.AppID,
		"&secret=", config.AppSecret,
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
	//c.JSON(http.StatusOK, gin.H{"openid": openid, "token": "123456"})
	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"openid": openid,
		"token":  "123456",
	})
}

// @Summary 获取支付订单
// @Tags 	微信公众号
// @Param   openid path string true "用户 openid"
// @Param   certid path string true "证书 id"
// @Success 200 {string} json "{"prepay_id":string}"
// @Router  /api/v1/weixin/wxorder/{openid}/{certid} [get]
func WXPayUnifyOrderReq(c *gin.Context) {
	appG := app.Gin{c}
	ip := c.ClientIP()
	openid := c.Param("openid")
/*
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
*/
	out_trade_no := UniqueId()
	fmt.Println(ip)
//	fmt.Println(result)

	client := wxpay.NewClient(wxpay.NewAccount(config.AppID, config.MchID, config.ApiKey, false))
	params := make(wxpay.Params)
//	params.SetString("body", result.CertName).
	params.SetString("body", "坤腾-证书").
		SetString("out_trade_no", out_trade_no).
	//	SetInt64("total_fee", int64(price)).
		SetInt64("total_fee", 1).
		SetString("spbill_create_ip", ip).
		SetString("notify_url", config.Notify_url).
		SetString("openid", openid).
		SetString("trade_type", "JSAPI")

	
	p, err := client.UnifiedOrder(params)
	if err != nil {
		appG.Response(http.StatusOK, false, e.INVALID_PARAMS, nil)
	}

	prepay_id := p.GetString("prepay_id")
	appid := p.GetString("appid")
	appG.Response(http.StatusOK, true, e.SUCCESS, map[string]string{
		"prepay_id": prepay_id,
		"appid":     appid,
		"price":  "1",
		"apikey":  config.ApiKey,
		"orderid": out_trade_no,
		"name": "test",
	})
	//c.JSON(http.StatusOK, gin.H{"prepay_id": prepay_id, "appid": appid})
}

type WXPayNotifyReq struct {
	Return_code            string `xml:"return_code"`
	Return_msg             string `xml:"return_msg"`
	Appid                  string `xml:"appid"`
	Mch_id                 string `xml:"mch_id"`
	Nonce                  string `xml:"nonce_str"`
	Sign                   string `xml:"sign"`
	Result_code            string `xml:"result_code"`
	Openid                 string `xml:"openid"`
	Is_subscribe           string `xml:"is_subscribe"`
	Trade_type             string `xml:"trade_type"`
	Bank_type              string `xml:"bank_type"`
	Total_fee              string `xml:"total_fee"`
	Fee_type               string `xml:"fee_type"`
	Cash_fee               int    `xml:"cash_fee"`
	Cash_fee_type          string `xml:"cash_fee_type"`
	Transaction_id         string `xml:"transaction_id"`
	Out_trade_no           string `xml:"out_trade_no"`
	Attach		       string `xml:"attach"`
	Time_end               string `xml:"time_end"`
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

func WXPayCallback(c *gin.Context) {
	var mr WXPayNotifyReq
	var ms WXPayNotifyResp
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
		ms.Return_code = "SUCCESS"
		ms.Return_msg = "OK"
	} else {
		ms.Return_code = "FAIL"
		ms.Return_msg = "failed to verify sign, please retry!"
	}

	c.XML(http.StatusOK, ms)
}

// @Summary 申请退款
// @Tags 	微信公众号
// @Param   out_trade_no path string true "付款订单号"
// @Success 200 {string} json "{"msg":string, "extra":}"
// @Router  /api/v1/weixin/wxrefund/{out_trade_no} [get]
func WXPayRefund(c *gin.Context) {
	out_trade_no := c.Param("out_trade_no")
	out_refund_no := UniqueId()
	//通过订单号去支付成功的数据库表中查找是否有此订单，并取出相应的total_fee,设置refund_fee

	account := wxpay.NewAccount(config.AppID, config.MchID, config.ApiKey, false)
	//	account.SetCertData(Weixin_cert)

	client := wxpay.NewClient(account)
	params := make(wxpay.Params)
	params.SetString("out_trade_no", out_trade_no).
		SetString("out_refund_no", out_refund_no).
		SetInt64("total_fee", 300).
		SetInt64("refund_fee", 300)

	p, err := client.Refund(params)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "UnifyOder failed"})
	}
	//更新退款成功的数据库表，记录退款成功状态
	c.JSON(http.StatusOK, gin.H{"msg": "SUCCESS", "extra": p})
}

// @Summary 查询退款
// @Tags 	微信公众号
// @Param   out_trade_no path string true "付款订单号"
// @Success 200 {string} json "{"msg":string, "extra":}"
// @Router  /api/v1/weixin/wxquery/{out_trade_no} [get]
func WXPayRefundQuery(c *gin.Context) {
	out_refund_no := c.Param("out_refund_no")

	client := wxpay.NewClient(wxpay.NewAccount(config.AppID, config.MchID, config.ApiKey, false))
	params := make(wxpay.Params)
	params.SetString("out_refund_no", out_refund_no)

	p, err := client.RefundQuery(params)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "UnifyOder failed"})
	}

	c.JSON(http.StatusOK, gin.H{"msg": "SUCCESS", "extra": p})
}

func wxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	signCalc := wxpayCalcSign(needVerifyM, config.ApiKey)

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
