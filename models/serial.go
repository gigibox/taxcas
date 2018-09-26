package models

import (
	"fmt"
	"strings"
)

const serialCol = "serial_log"

var areaCode = map[string]string{
	"北京": "01",
	"天津": "02",
	"河北": "03",
	"山西": "04",
	"内蒙古": "05",
	"辽宁": "06",
	"吉林": "07",
	"黑龙江": "08",
	"上海": "09",
	"江苏": "10",
	"浙江": "11",
	"安徽": "12",
	"福建": "13",
	"江西": "14",
	"山东": "15",
	"河南": "16",
	"湖北": "17",
	"湖南": "18",
	"广东": "19",
	"广西": "20",
	"海南": "21",
	"重庆": "22",
	"四川": "23",
	"贵州": "24",
	"云南": "25",
	"西藏": "26",
	"陕西": "27",
	"甘肃": "28",
	"青海": "29",
	"宁夏": "30",
	"新疆": "31",
	"台湾": "32",
	"香港": "33",
	"澳门": "34",
}

// 证书编号生成规则 [日期] + [地区编号] + [6位顺序号] , (e.g.): 20180901000001

type serialTables struct {
	Prefix string
	Count int

}

func GenerateCertSN(studyDate, province, certid string) (string, bool) {
	var result string

	d := strings.Split(studyDate, "-")
	serialPrefix := d[0] + d[1] + areaCode[province]

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