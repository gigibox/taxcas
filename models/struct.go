package models

const (
	NotPaid = iota
	Paid
	Pending
	Verifying
	Reject
	Passed
	Refund
	Refunded
	Error
)

var StatusMsg = map[int]string{
	NotPaid:	"未支付",
	Paid:		"已支付",
	Pending:	"待审核",
	Verifying:	"审核中",
	Reject:		"已拒绝",
	Passed:		"已通过",
	Refund:		"退款中",
	Refunded:	"已退款",
	Error:		"错误状态",
}

var ActionMsg = map[string]int {
	"export" : Pending,
	"verify" : Verifying,
	"passed" : Passed,
	"reject" : Reject,
	"error"  : Error,
}

type Coord struct {
	Str string
	Font string			`json:"font" form:"font"`
	FontSize float64	`json:"font_size" form:"font_size"`
	X int
	Y int
}

// 绘图位置
type ImageDesigner struct {
	ImgName string		`json:"img_name" form:"img_name"`
	Name Coord			`json:"name" form:"name"`
	EnglishName Coord	`json:"english_name" form:"english_name"`
	SerialNumber Coord	`json:"serial_number" form:"serial_number"`
	PersonalID Coord	`json:"personal_id" form:"personal_id"`
	Date Coord			`json:"date" form:"date"`
}

type User struct {
	WechatID 	string	`json:"wechat_id" valid:"Required; MaxSize(60)"` // 微信ID
	Name 		string	`json:"name" valid:"Required; MaxSize(60)"`// 姓名
	EnglishName string	`json:"english_name" valid:"Required; MaxSize(60)"`// 拼音
	Phone		int		`json:"phone" valid:"Required"` // 电话号码
	PersonalID 	string	`json:"personal_id" valid:"Required"` // 身份证号
	Job 		string	`json:"job" valid:"Required"` // 职位
	Address 	string	`json:"address"` // 地址
	Province	string	`json:"province" valid:"Required; MaxSize(60)"`
	City		string	`json:"city" valid:"Required; MaxSize(60)"`
	District	string	`json:"district" valid:"Required; MaxSize(60)"`
	Company 	string	`json:"company_name" valid:"Required"` // 公司名称
	CompanyNature string `json:"company_nature" valid:"Required"`
	CompanyScale string	`json:"company_scale" valid:"Required"`
}

type Applicant struct {
	User
	CertID			string	`json:"cert_id"`
	CertName	 	string	`json:"cert_name" valid:"Required; MaxSize(60)"`	// 证书名称
	StudyDate		string	`json:"study_date" valid:"Required; MaxSize(60)"`
	CertType		int		`json:"cert_type"`
}

// 证书类型表
type C_certs struct {
	CertID		 int		`json:"cert_id"`// 证书索引
	CertName	 string		`json:"cert_name" form:"cert_name" valid:"Required"`	// 证书名称
	Authority	 string		`json:"authority" form:"authority"`	// 证书颁发机构
	Price		 int	    `json:"price" form:"price"`// 申请费用
	Status		 string		`json:"status" form:"status" valid:"Required"`	// 申请状态
	ImageDesign  ImageDesigner `json:"image_design"`// 内容位置
}

// 管理员
type C_admin struct {
	Username string `json:"username" valid:"Required; MaxSize(50)"`
	Password string `json:"password" valid:"Required; MaxSize(50)"`
}

// 用户表
type C_users struct {
	User
	Certs       []map[string]int
}

// 证书申领流水
type C_Apply struct {
	Applicant
	PayAmount       int		`json:"pay_amount"`		// 支付金额
	PayOrder		int		`json:"pay_order"`		// 支付订单
	PayStatus		int		`json:"pay_status"`		// 支付状态
	SerialNumber    string	`json:"serial_number"`  // 证书编号
	ImageSaveUrl	string	`json:"image_save_url"`
	ApplyStatus     int		`json:"apply_status"`	// 申请状态
	ApplyStatusMsg  string	`json:"apply_status_msg"`//申请状态信息
	ApplyDate		int  `json:"apply_date"`		// 申请时间
}