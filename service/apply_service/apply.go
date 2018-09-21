package apply_service

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"taxcas/models"
	"taxcas/pkg/export"
	"taxcas/pkg/logging"
	"taxcas/service/user_service"
	"time"
)

type S_Apply struct {
	Collection string
	Data models.C_Apply
}

func (this *S_Apply) CheckApplyExist() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("applicant.user.wechatid", this.Data.WechatID, this.Collection, &result)
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

func (this *S_Apply) Add() (bool, error) {
	return models.MgoInsert(this.Data, this.Collection)
}

func New(col string, commit models.Applicant) (*S_Apply) {
	return &S_Apply{
		Collection: col,
		Data: models.C_Apply{
			Applicant: commit,
			ApplyDate: int(time.Now().Unix()),
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

func GetApplyList(id, action string, page, limit int) (interface{}) {
	result := []models.C_Apply{}

	if page > 0 {
		page -= 1
	}

	key, val := parseAction(action)

	// 统计符合条件的总数
	count, _ := models.MgoCountQuery(key, val, "cert" + id + "_apply")

	// 查询
	models.MgoFind(key, val, "cert" + id + "_apply", page, limit, &result)

	return map[string]interface{} {
		"count": count,
		"list" : result,
	}
}

// 根据请求类型判断需要导出的字段
var title = map[string][]string {
	"export" : []string{"申请证书", "微信ID", "申请人", "身份证号", "申请时间", "支付金额"},
}
func ExportFile(certid, act string) (string, error) {
	code, ok := models.ActionMsg[act]
	if !ok {
		return "", nil
	}

	// 查询结果
	docs := []models.C_Apply{}
	key, val := parseAction(act)
	models.MgoFind(key, val, "cert" + certid + "_apply", 0, 0, &docs)

	// 文件名已查询条件 + 时间 命名
	timestamp := strconv.Itoa(int(time.Now().Unix()))
	filename := models.StatusMsg[code] + "-" +  timestamp + ".csv"

	// 创建 csv 文件
	f, err := os.Create(export.GetExportFullPath() + filename)
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
		t := time.Unix(int64 (docs[i].ApplyDate), 0)

		row := []string{
			docs[i].CertName,
			docs[i].WechatID,
			docs[i].Name,
			docs[i].PersonalID,
			t.Format("2018-09-19 17:27:27"),
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
		ApplyStatus: newCode,
		ApplyStatusMsg: newMsg,
	}

	models.MgoUpdateAll(key, val, "cert" + certid + "_apply", apply)

	// 更新用户表状态
	for i := range docs {
		user := models.User{
			WechatID : docs[i].WechatID,
		}
		user_service.UpdateCerts(user, certid, newCode)
	}

	return filename, err
}

func UpdateApplicantsByFile(certid, act, file string) (int, int) {
	var succeed, failure int

	statusCode, ok := models.ActionMsg[act]
	if !ok {
		return succeed, failure
	}

	statusMsg := models.StatusMsg[statusCode]

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
			fmt.Println(err)
			failure ++
			continue
		}
		// 状态
		wxid  := record[1]

		doc := models.C_Apply{
			ApplyStatus: statusCode,
			ApplyStatusMsg: statusMsg,
		}

		// 根据微信id 更新申请订单状态
		if ok , err := models.MgoUpsert("applicant.user.wechatid", wxid, "cert" + certid + "_apply", doc); !ok {
			logging.Warn("UpdateApplicantsByFile ERROR:", err)
			failure ++
		} else {
			succeed ++

			// 判断为已通过, 生成证书编号
			if statusCode == models.Passed {

			}

			// 判断为退款请求, 发起退款申请
			if statusCode == models.Refunded {

			}

			// 修改用户申请状态
			user := models.User{
				WechatID : wxid,
			}
			user_service.UpdateCerts(user, certid, statusCode)

			// 推送微信提醒
		}
	}

	return succeed, failure
}