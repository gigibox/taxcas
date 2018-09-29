package apply_service

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"taxcas/models"
	"taxcas/pkg/export"
	"taxcas/pkg/logging"
	"taxcas/service/user_service"
	"taxcas/service/weixin_service"
	"time"
)

type S_Apply struct {
	Collection string
	Data       models.C_Apply
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

	if result.Status == "enable" {
		return true, nil
	}

	return false, nil
}

func (this *S_Apply) UpdateSerialNumber() bool {
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

func (this *S_Apply) UpdateStatus() bool {
	statusCode := this.Data.ApplyStatus

	// 根据身份证号 更新申请订单状态
	if this.Data.PersonalID != "" {
		if ok, err := models.MgoUpdate("applicant.user.personalid", this.Data.PersonalID, this.Collection, this.Data); !ok {
			logging.Warn("Update applicant status:", err)
			return false
		}
	} else if this.Data.WechatID != "" {
		if ok, err := models.MgoUpdate("applicant.user.wechatid", this.Data.WechatID, this.Collection, this.Data); !ok {
			logging.Warn("Update applicant status:", err)
			return false
		}
	} else {
		return false
	}

	// 判断为退款请求, 发起退款申请
	if statusCode == models.Refunded {
		if ok, err := weixin_service.WXPayRefund(this.Data.PayOrder); !ok {
			logging.Error("订单%s, 退款失败. %s", this.Data.PayOrder, err)
		}

	}

	// 修改用户申请状态
	user := models.User{
		PersonalID: this.Data.PersonalID,
	}
	user_service.UpdateCerts(user, this.Data.CertID, this.Data.ApplyStatus)

	// 推送微信提醒

	return true
}

func New(col string, commit models.Applicant) *S_Apply {
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

func GetApplyList(id, action string, page, limit int) interface{} {
	result := []models.C_Apply{}

	if page > 0 {
		page -= 1
	}

	key, val := parseAction(action)

	// 统计符合条件的总数
	count, _ := models.MgoCountQuery(key, val, "cert"+id+"_apply")

	// 查询
	models.MgoFind(key, val, "cert"+id+"_apply", page, limit, &result)

	return map[string]interface{}{
		"count": count,
		"list":  result,
	}
}

func GetApplyByPID(certid, pid string, doc *models.C_Apply) (bool, error) {
	return models.MgoFindOne("applicant.user.personalid", pid, "cert"+certid+"_apply", doc)
}

func GetApplyByOpenid(certid, openid string, doc *models.C_Apply) (bool, error) {
	return models.MgoFindOne("applicant.user.wechatid", openid, "cert"+certid+"_apply", doc)
}

func GetApplyBySN(certid, sn string, doc *models.C_Apply) (bool, error) {
	return models.MgoFindOne("serialnumber", sn, "cert"+certid+"_apply", doc)
}

// 根据请求类型判断需要导出的字段
var title = map[string][]string{
	"export": []string{"编号", "申请证书", "申请人", "身份证号", "申请时间", "支付金额"},
}

func ExportFile(certid, act string) (string, error) {
	code, ok := models.ActionMsg[act]
	if !ok {
		return "", nil
	}

	// 查询结果
	docs := []models.C_Apply{}
	key, val := parseAction(act)
	models.MgoFind(key, val, "cert"+certid+"_apply", 0, 0, &docs)

	// 文件名已查询条件 + 时间 命名
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	filecsv := export.GetExportExcelPath() + models.StatusMsg[code] + "-" + timestamp + ".csv"

	// 创建 csv 文件
	f, err := os.Create(export.GetRuntimePath() + filecsv)
	if err != nil {
		panic(err)
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
			docs[i].SerialNumber,
			docs[i].CertName,
			docs[i].Name,
			docs[i].PersonalID,
			docs[i].StudyDate,
			strconv.Itoa(docs[i].PayAmount),
		}

		data = append(data, row)
	}

	w.WriteAll(data)

	// 更新申请状态
	var newCode int

	if code == models.Pending {
		newCode = models.Verifying
	} else if code == models.Reject {
		newCode = models.Refunding
	}

	newMsg := models.StatusMsg[newCode]

	// 更新申请订单表
	apply := models.C_Apply{
		ApplyStatus:    newCode,
		ApplyStatusMsg: newMsg,
	}

	// 一次更新所有状态
	models.MgoUpdateAll(key, val, "cert"+certid+"_apply", apply)

	// 更新用户表状态
	for i := range docs {
		user := models.User{
			PersonalID: docs[i].PersonalID,
		}
		user_service.UpdateCerts(user, certid, newCode)
	}

	return filecsv, err
}

func UpdateApplicants(certid, act, file string, pids []string) (int, int) {
	var succeed, failure int

	statusCode, ok := models.ActionMsg[act]
	if !ok {
		return succeed, failure
	}

	statusMsg := models.StatusMsg[statusCode]

	applyService := S_Apply{
		Collection: "cert" + certid + "_apply",
		Data: models.C_Apply{
			ApplyStatus:    statusCode,
			ApplyStatusMsg: statusMsg,
		},
	}

	// 根据导入数据处理
	if file != "" {
		// 解析cav文件
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		reader := csv.NewReader(f)

		// 跳过第一行
		reader.Read()

		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			} else if nil != err {
				logging.Error(err)
				failure++
				continue
			}

			applyService.Data.PersonalID = record[3]
			if ok := applyService.UpdateStatus(); ok {
				succeed++
			} else {
				failure++
			}
		}
	} else {
		// 手动选择
		for i := range pids {
			applyService.Data.PersonalID = pids[i]
			if ok := applyService.UpdateStatus(); ok {
				succeed++
			} else {
				failure++
			}
		}

	}

	return succeed, failure
}
