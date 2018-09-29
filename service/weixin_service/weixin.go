package weixin_service

import (
	"taxcas/models"
	"github.com/objcoding/wxpay"
	"crypto/md5"
	"encoding/base64"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"
	"encoding/hex"
	"io/ioutil"
	"io"
	"errors"
	"crypto/rand"
	"fmt"
	"os"
	"time"
	"encoding/json"
	"net/http"
	"strings"
	"bytes"
)

const (
	Col_order  = "order"
	Col_refund = "refund"
)

var accessToken *AccessToken
func Add(order models.WXPayNotifyReq, col string) (bool, error) {
	return models.MgoInsert(order, col)
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

func WXPayRefund(out_trade_no string) (bool, error) {
        out_refund_no := UniqueId()
        //通过订单号去支付成功的数据库表中查找是否有此订单，并取出相应的total_fee,设置refund_fee

        account := wxpay.NewAccount(setting.WeixinSetting.AppID, setting.WeixinSetting.MchID, setting.WeixinSetting.ApiKey, false)
        account.SetCertData(upload.GetApiCertFullPath())

        result := models.WXPayNotifyReq{}
        isExist, err := models.MgoFindOne("out_trade_no", out_trade_no, Col_order, &result)
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

type AccessToken struct {
        Token      string `json:"access_token"`
        ExpiresIn  int    `json:"expires_in"`
        CreateTime int64  
}


func postReq(url, reqBody string) (string, error) {
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
}

/**
* @note 获取微信的access_token
 * @params *accessToken 用来存放token的结构体
  * @return token string
*/
func getAccessToken() string {
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
}

