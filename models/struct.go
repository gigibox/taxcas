package models

const (
	NotPaid = iota
	Paid
	Pending
	Verifying
	Reject
	Passed
	Refunding
	Refunded
	Error
)

var StatusMsg = map[int]string{
	NotPaid:   "未支付",
	Paid:      "已支付",
	Pending:   "待审核",
	Verifying: "审核中",
	Reject:    "已拒绝",
	Passed:    "已通过",
	Refunding: "退款中",
	Refunded:  "已退款",
	Error:     "错误状态",
}

var ActionMsg = map[string]int{
	"export":   Pending,
	"verify":   Verifying,
	"passed":   Passed,
	"reject":   Reject,
	"refunded": Refunded,
}

type Coord struct {
	Str      string
	Font     string  `json:"font" form:"font"`
	FontSize float64 `json:"font_size" form:"font_size"`
	X        int
	Y        int
}

// 绘图位置
type ImageDesigner struct {
	ImgName      string `json:"img_name" form:"img_name"`
	Name         Coord  `json:"name" form:"name"`
	EnglishName  Coord  `json:"english_name" form:"english_name"`
	SerialNumber Coord  `json:"serial_number" form:"serial_number"`
	PersonalID   Coord  `json:"personal_id" form:"personal_id"`
	Date         Coord  `json:"date" form:"date"`
}

type User struct {
	WechatID      string `json:"wechat_id" valid:"Required; MaxSize(60)"`    // 微信ID
	Name          string `json:"name" valid:"Required; MaxSize(60)"`         // 姓名
	EnglishName   string `json:"english_name" valid:"Required; MaxSize(60)"` // 拼音
	Phone         int    `json:"phone" valid:"Required"`                     // 电话号码
	PersonalID    string `json:"personal_id" valid:"Required"`               // 身份证号
	Job           string `json:"job" valid:"Required"`                       // 职位
	Address       string `json:"address"`                                    // 地址
	Province      string `json:"province" valid:"Required; MaxSize(60)"`
	City          string `json:"city" valid:"Required; MaxSize(60)"`
	District      string `json:"district" valid:"Required; MaxSize(60)"`
	Company       string `json:"company_name" valid:"Required"` // 公司名称
	CompanyNature string `json:"company_nature" valid:"Required"`
	CompanyScale  string `json:"company_scale" valid:"Required"`
}

type Applicant struct {
	User
	CertID    string `json:"cert_id"`
	CertName  string `json:"cert_name" valid:"Required; MaxSize(60)"` // 证书名称
	StudyDate string `json:"study_date" valid:"Required; MaxSize(60)"`
	CertType  int    `json:"cert_type"`
}

// 证书类型表
type C_certs struct {
	CertID      string        `json:"cert_id"`                                     // 证书索引
	CertName    string        `json:"cert_name" form:"cert_name" valid:"Required"` // 证书名称
	Authority   string        `json:"authority" form:"authority"`                  // 证书颁发机构
	Price       int           `json:"price" form:"price"`                          // 申请费用
	Status      string        `json:"status" form:"status" valid:"Required"`       // 申请状态
	ImageDesign ImageDesigner `json:"image_design"`                                // 内容位置
}

// 管理员
type C_admin struct {
	Username string `json:"username" valid:"Required; MaxSize(50)"`
	Password string `json:"password" valid:"Required; MaxSize(50)"`
}

// 用户表
type C_users struct {
	User
	Certs []map[string]int
}

// 证书申领流水
type C_Apply struct {
	Applicant
	PayAmount      int    `json:"pay_amount"`    // 支付金额
	PayOrder       int    `json:"pay_order"`     // 支付订单
	PayStatus      int    `json:"pay_status"`    // 支付状态
	SerialNumber   string `json:"serial_number"` // 证书编号
	ImageSaveUrl   string `json:"image_save_url"`
	ApplyStatus    int    `json:"apply_status"`     // 申请状态
	ApplyStatusMsg string `json:"apply_status_msg"` //申请状态信息
	ApplyDate      int    `json:"apply_date"`       // 申请时间
}

var GlobalDesigner = ImageDesigner{
	Name: Coord{
		Str:      "李雷",
		Font:     "微软雅黑",
		FontSize: 20,
		X:        230,
		Y:        240,
	},
	EnglishName: Coord{
		Str:      "LiLei",
		Font:     "微软雅黑",
		FontSize: 12,
		X:        430,
		Y:        240,
	},
	PersonalID: Coord{
		Str:      "110010201010010101",
		Font:     "微软雅黑",
		FontSize: 12,
		X:        530,
		Y:        140,
	},
	SerialNumber: Coord{
		Str:      "20180914012345",
		Font:     "微软雅黑",
		FontSize: 12,
		X:        830,
		Y:        240,
	},
	Date: Coord{
		Str:      "2018 年 9 月 14 日",
		Font:     "微软雅黑",
		FontSize: 12,
		X:        930,
		Y:        240,
	},
}

type WXPayNotifyReq struct {
	Return_code    string `xml:"return_code" json:"return_code"`
	Return_msg     string `xml:"return_msg" json:"return_msg"`
	Appid          string `xml:"appid" json:"appid"`
	Mch_id         string `xml:"mch_id" json:"mch_id"`
	Nonce          string `xml:"nonce_str" json:"nonce_str"`
	Sign           string `xml:"sign" json:"sign"`
	Result_code    string `xml:"result_code" json:"result_code"`
	Openid         string `xml:"openid" json:"openid"`
	Is_subscribe   string `xml:"is_subscribe" json:"is_subscribe"`
	Trade_type     string `xml:"trade_type" json:"trade_type"`
	Bank_type      string `xml:"bank_type" json:"bank_type"`
	Total_fee      string `xml:"total_fee" json:"total_fee"`
	Fee_type       string `xml:"fee_type" json:"fee_type"`
	Cash_fee       int    `xml:"cash_fee" json:"cash_fee"`
	Cash_fee_type  string `xml:"cash_fee_type" json:"cash_fee_type"`
	Transaction_id string `xml:"transaction_id" json:"transaction_id"`
	Out_trade_no   string `xml:"out_trade_no" json:"out_trade_no"`
	Attach         string `xml:"attach" json:"attach"`
	Time_end       string `xml:"time_end" json:"time_end"`
}

type WXPayNotifyResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}
