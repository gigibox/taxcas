package models

import (
	"fmt"
	"strings"
)

const serialCol = "serial_log"

var areaCode = map[string]string{
	"北京市": "01",
	"天津市": "02",
	"河北省": "03",
	"山西省": "04",
	"内蒙古自治区": "05",
	"辽宁省": "06",
	"吉林省": "07",
	"黑龙江省": "08",
	"上海市": "09",
	"江苏省": "10",
	"浙江省": "11",
	"安徽省": "12",
	"福建省": "13",
	"江西省": "14",
	"山东省": "15",
	"河南省": "16",
	"湖北省": "17",
	"湖南省": "18",
	"广东省": "19",
	"广西壮族自治区": "20",
	"海南省": "21",
	"重庆市": "22",
	"四川省": "23",
	"贵州省": "24",
	"云南省": "25",
	"西藏自治区": "26",
	"陕西省": "27",
	"甘肃省": "28",
	"青海省": "29",
	"宁夏回族自治区": "30",
	"新疆维吾尔自治区": "31",
	"台湾省": "32",
	"香港特别行政区": "33",
	"澳门特别行政区": "34",
}

// 证书编号生成规则 [日期] + [地区编号] + [6位顺序号] , (e.g.): 20180901000001

type serialTables struct {
	Prefix string
	Count int

}

func GenerateCertSN(studyDate, province, certid string) (string, bool) {
	var result string

	ac := areaCode[province]
	if ac == "" {
		ac = "99"
	}

	d := strings.Split(studyDate, "-")
	serialPrefix := d[0] + d[1] + ac

	doc := serialTables{
		Prefix : serialPrefix + "_cert" + certid,
		Count : 0,
	}

	// 获取顺序号
	MgoFindOne("prefix", doc.Prefix, serialCol, &doc)

	doc.Count += 1

	// 更新记录
	if _, err := MgoUpsert("prefix", doc.Prefix, serialCol, doc); err != nil {
		return result, false
	}

	result = fmt.Sprintf("%s%06d", serialPrefix, doc.Count)

	return result, true
}