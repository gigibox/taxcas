package apply_service

import (
	"encoding/csv"
	"gopkg.in/mgo.v2/bson"
	"io"
	"os"
	"strconv"
	"taxcas/models"
	"taxcas/pkg/export"
	"taxcas/pkg/logging"
	"taxcas/pkg/util"
	"taxcas/service/msg_service"
	"taxcas/service/user_service"
	"taxcas/service/weixin_service"
	"time"
)

type S_Apply struct {
	Collection string
	Data models.C_Apply
}

func (this *S_Apply) CheckApplyExistByWX() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("applicant.user.wechatid", this.Data.WechatID, this.Collection, &result)
}

func (this *S_Apply) CheckApplyExistByID() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("applicant.user.personalid", this.Data.PersonalID, this.Collection, &result)
}

func (this *S_Apply) CheckCertByName() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certname", this.Data.CertName, "certs", &result)
}

func (this *S_Apply) CheckApplyStatus() (bool, error) {
	result := models.C_certs{}
	models.MgoCheckKeyExist("certname", this.Data.CertName, "certs", &result)

	if result.Status == "enabled" {
		return true, nil
	}

	return false, nil
}

func (this *S_Apply) UpdateSerialNumber() (bool){
	sn, ok := models.GenerateCertSN(this.Data.StudyDate, this.Data.Province, this.Data.CertID)
	if !ok {
		return false
	} else {
		this.Data.SerialNumber = sn
		return true
	}
}

func (this *S_Apply) Add() (bool, error) {
	return models.MgoInsert(this.Data, this.Collection)
}

func (this *S_Apply) Update() (bool, error) {
	return models.MgoUpdate("applicant.user.personalid", this.Data.PersonalID, this.Collection, this.Data);
}

func (this *S_Apply) UpdateStatus() (bool) {
		// 根据身份证号 更新申请订单状态
		if ok , err := models.MgoUpdate("applicant.user.personalid", this.Data.PersonalID, this.Collection, this.Data); !ok {
			logging.Debug("Update applicant status:", err)
			return false
		}

		// 修改用户申请状态
		user := models.User{
			WechatID : this.Data.WechatID,
		}
		user_service.UpdateCerts(user, this.Data.CertID, this.Data.ApplyStatus)

		return true
}

func New(col string, commit models.Applicant) (*S_Apply) {
	return &S_Apply{
		Collection: col,
		Data: models.C_Apply{
			Applicant: commit,
			ApplyDate: time.Now().Unix(),
		},
	}
}

func parseAction(action string) (string, interface{}) {
	flag, ok := models.ActionMsg[action]
	if ok {
		return "applystatus", flag
	}

	return "", nil
}

func GetApplyList(id, action string, page, limit int, ext string) (interface{}) {
	doc := []models.C_Apply{}

	if page > 0 {
		page -= 1
	}

	key, val := parseAction(action)

	selecter :=	bson.M{key: val}

	// 输入框查询, 18位为身份证号, 否则查询姓名
	if ext != "" {
		if len(ext) == 18 {
			selecter = bson.M{key: val, "applicant.user.personalid": ext}
		} else {
			selecter = bson.M{key: val, "applicant.user.name": ext}
		}
	}

	// 统计符合条件的总数
	count, _ := models.MgoCountQuery(selecter, "cert" + id + "_apply")

	// 查询
	models.MgoFind(selecter, "cert" + id + "_apply", page, limit, &doc)

	return map[string]interface{} {
		"count": count,
		"list" : doc,
	}
}

func GetApplyByPID(certid, pid string, doc *models.C_Apply) (bool, error) {
	return models.MgoFindOne("applicant.user.personalid", pid, "cert" + certid + "_apply", doc)
}

func GetApplyByOpenid(certid, openid string, doc *models.C_Apply) (bool, error) {
	return models.MgoFindOne("applicant.user.wechatid", openid, "cert" + certid + "_apply", doc)
}

func GetApplyBySN(certid, sn string, doc *models.C_Apply) (bool, error) {
	return models.MgoFindOne("serialnumber", sn, "cert" + certid + "_apply", doc)
}

// 根据请求类型判断需要导出的字段
var title = map[string][]string {
	"export" : []string{"编号", "申请证书", "申请人", "身份证号", "申请时间", "支付金额"},
}
func ExportFile(certid, act string) (string, error) {
	code, ok := models.ActionMsg[act]
	if !ok {
		return "", nil
	}

	// 查询结果
	key, val := parseAction(act)

	selecter := bson.M{key: val}
	if val == models.Reject {
		selecter = bson.M{key: val, "paystatus": models.Paid}
	}

	docs := []models.C_Apply{}
	models.MgoFind(selecter, "cert" + certid + "_apply", 0, 0, &docs)

	if len(docs) == 0 {
		logging.Debug("Export csv File, condition not queried ")
		return "none", nil
	}

	// 文件名: 证书名称-查询条件-日期.csv
	tm := time.Unix(time.Now().Unix(), 0)
	filecsv := export.GetExportExcelPath() + docs[0].CertName + "-" + models.StatusMsg[code] + "-" +  tm.Format("20060102") + ".csv"

	// 创建 csv 文件
	f, err := os.Create(export.GetRuntimePath() + filecsv)
	if err != nil {
		logging.Error("Creat csv file fail: ",err)
		return "", err
	}
	defer f.Close()

	// UTF-8 BOM
	f.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(f)

	// 写入头信息
	w.Write(title[act])

	data := [][]string{}

	for i := range docs {
		row := []string{
			docs[i].SerialNumber + "\t", // 避免数字以科学计数法显示
			docs[i].CertName,
			docs[i].Name,
			docs[i].PersonalID + "\t",
			docs[i].StudyDate,
			strconv.Itoa(docs[i].PayAmount),
		}

		data = append(data, row)
	}

	w.WriteAll(data)

	// 更新申请状态
	var newCode int
	var newMsg  string
	var newData bson.M

	if code == models.Pending {
		newCode = models.Verifying
		newMsg = models.StatusMsg[newCode]

		// 更新申请订单表
		newData = bson.M{"$set": bson.M{"applystatus": newCode, "applystatusmsg": newMsg}}
	} else if code == models.Reject {
		newCode = models.Refunding
		newMsg = models.StatusMsg[newCode]
		newTime := time.Now().Unix()

		// 更新支付状态
		newData = bson.M{"$set": bson.M{"paystatus": newCode, "applystatusmsg": newMsg, "paytime": newTime}}
	}

	// 一次更新所有状态
	if ok, err := models.MgoUpdateAll(selecter, "cert" + certid + "_apply", newData); !ok {
		logging.Error("Update all apply status failure: ", err)
	}

	// 更新用户表状态
	for i := range docs {
		user := models.User{
			WechatID : docs[i].WechatID,
		}
		user_service.UpdateCerts(user, certid, newCode)
	}

	return filecsv, err
}

func UpdateApplicants(certid, act, file string, pids []string) (int, int) {
	var succeed, failure int

	// 手动拒绝参数为身份证号数组
	pidArry := pids

	statusCode, ok := models.ActionMsg[act]
	if !ok {
		return succeed, failure
	}

	statusMsg := models.StatusMsg[statusCode]

	applyService := S_Apply{
		Collection: "cert" + certid + "_apply",
	}

	// 解析csv文件
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			logging.Error(err)
			return 0, 0
		}
		defer f.Close()

		reader := csv.NewReader(f)

		// 跳过第一行
		reader.Read()

		pidArry = []string{}
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if nil != err {
				logging.Error(err)
				failure ++
				continue
			}

			// 去除空格和制表符, 读取身份证号
			pidArry = append(pidArry, util.CompressStr(record[3]))
		}
	}

	for i := range pidArry{
		if ext, _ := GetApplyByPID(certid, pidArry[i], &applyService.Data); !ext {
			failure ++
			continue
		}

		// 避免重复导入
		if applyService.Data.ApplyStatusMsg != statusMsg {
			applyService.Data.ApplyStatusMsg = statusMsg

			// 如果是已拒绝状态, 申请状态不改变
			if applyService.Data.ApplyStatus != models.Reject {
				applyService.Data.ApplyStatus = statusCode
			}

			// 更新退款状态
			if statusCode == models.Refunded {
				applyService.Data.PayTime = time.Now().Unix()
				applyService.Data.PayStatus = statusCode
			}

			// 判断为退款请求, 发起退款申请
			if statusCode == models.Refunded && applyService.Data.PayOrder != "" {
				if ok, err := weixin_service.WXPayRefund(applyService.Data.PayOrder); !ok {
					logging.Error("Pay order: %s, Refund failure: %s", applyService.Data.PayOrder, err)
					failure ++
					continue
				}
			}

			// 推送微信提醒
			msg_service.Send(statusCode, &applyService.Data)

			if ok := applyService.UpdateStatus(); ok {
				succeed ++
			} else {
				failure ++
			}

		}
	}

	return succeed, failure
}
