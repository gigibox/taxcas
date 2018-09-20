package cert_service

import "C"
import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"taxcas/models"
	"taxcas/pkg/export"
	"time"
)

type S_cert struct {
	Collection string
	Data models.C_certs
}

const (
	col_certs = "certs"
)

func (this *S_cert) CheckExist() (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certname", this.Data.CertName, this.Collection, &result)
}

func (this *S_cert) Add() (bool, error) {
	var err error

	this.Data.Status = "enable"
	this.Data.CertID, err = models.MgoCountCollection(this.Collection)

	// 证书模板设计, 暂时写死
	this.Data.ImageDesign.Name = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.EnglishName = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.PersonalID = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.SerialNumber = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}
	this.Data.ImageDesign.Date = models.Coord{
		Font: "微软雅黑",
		FontSize: 12,
		X: 230,
		Y: 240,
	}

	// 序号从1开始
	this.Data.CertID++
	if err != nil {
		log.Println(err)
		return false, err
	}

	return models.MgoInsert(this.Data, this.Collection)
}

func GetAllCertName() (interface{}) {
	type simpleCert struct {
		Id string `json:"cert_id"`
		Name string `json:"cert_name"`
		Status string `json:"status"`
	}

	certs := []simpleCert{}
	result := []models.C_certs{}

	models.MgoFindAll(col_certs, &result)
	for i, _ := range result {
		t := simpleCert{Id:strconv.Itoa(result[i].CertID), Name:result[i].CertName, Status:result[i].Status}
		certs = append(certs, t)
	}

	return certs
}

func CheckExistByID(id int) (bool, error) {
	result := models.C_certs{}
	return models.MgoCheckKeyExist("certid", id, col_certs, &result)
}

func SignCertImage(design models.ImageDesigner) (string, error) {
	imageSaveUrl, err := models.SignImage(design)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return imageSaveUrl, nil
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

var title = map[string][]string {
				"export" : []string{"申请证书", "微信ID", "申请人", "身份证号", "申请时间", "支付金额", "申请状态"},
			}

func ExportFile(id, act string) (string, error) {
	code, ok := models.ActionMsg[act]
	if !ok {
		return "", nil
	}

	// 查询结果
	docs := []models.C_Apply{}
	key, val := parseAction(act)
	models.MgoFind(key, val, "cert" + id + "_apply", 0, 0, &docs)

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

		row := []string{docs[i].CertName,
						docs[i].WechatID,
						docs[i].Name,
						docs[i].PersonalID,
						t.Format("2018-09-19 17:27:27"),
						strconv.Itoa(docs[i].PayAmount),
						models.StatusMsg[code]}

		data = append(data, row)
	}

	w.WriteAll(data)

	return filename, err
}

func parseStatusMsg(msg string) (int, bool) {
	for k, v := range models.StatusMsg{
		if v == msg {
			return k, true
		}
	}

	return 0, false
}

func ImportResult(id, act, file string) (int, int) {
	var succeed, failure int

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
		wechatID  := record[1]
		statusMsg := record[len(record) - 1]
		statusCode, ok := parseStatusMsg(statusMsg)
		if !ok {
			failure ++
			continue
		}

		doc := models.C_Apply{
			ApplyStatus: statusCode,
			ApplyStatusMsg: statusMsg,
		}

		// 根据微信id 更新申请订单状态
		if ok , _ := models.MgoUpsert("applicant.user.wechatid", wechatID, "cert" + id + "_apply", doc); !ok {
			failure ++
		} else {
			succeed ++
		}

		// 判断为已审核内容, 生成证书编号
	}

	return succeed, failure
}