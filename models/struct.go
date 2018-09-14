package models

// 管理员
type t_admin struct {
	Username string
	Password string
}

type coord struct {
	X int
	Y int
}

// 绘图位置
type position struct {
	Name coord
	EnglishName coord
	SerialNumber coord
	PersonalID coord
	Date coord
}

type C_certs struct {
	CertID		 int		// 证书索引
	CertName	 string		`valid:"Required; MaxSize(60)"`	// 证书名称
	CertPrefix	 string		`valid:"Required; MaxSize(50)"`	// 证书编号前缀
	Cost		 int	    // 申请费用
	Status		 string		`valid:"Required; MaxSize(30)"`	// 申请状态
	ImgName		 string		`valid:"Required; MaxSize(50)"`	// 图片路径
	Positions	 position	`valid:"Required;"`// 内容位置
}

type C_users struct {
	Name 		string	// 姓名
	EnglishName string	// 拼音
	Phone		int		// 电话号码
	PersonalID 	string	// 身份证号
	WechatID 	string	// 微信ID
	Company 	string	// 公司名称
	Job 		string	// 职位
	Address 	string	// 地址
}