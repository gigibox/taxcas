package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_EXIST_CERT:                "该证书名称已存在",
	ERROR_EXIST_CERT_FAIL:           "获取已存在证书失败",
	ERROR_NOT_EXIST_CERT:            "该证书不存在",
	ERROR_GET_CERTS_FAIL:            "获取所有证书失败",
	ERROR_COUNT_CERT_FAIL:           "统计证书失败",
	ERROR_ADD_CERT_FAIL:             "新增证书失败",
	ERROR_EDIT_CERT_FAIL:            "修改证书失败",
	ERROR_DELETE_CERT_FAIL:          "删除证书失败",
	ERROR_EXPORT_FILE_FAIL:          "导出文件失败",
	ERROR_IMPORT_FILE_FAIL:          "导入文件失败",
	ERROR_NOT_EXIST_USER:            "该用户不存在",
	ERROR_ADD_USER_FAIL:             "新增用户失败",
	ERROR_DELETE_USER_FAIL:          "删除用户失败",
	ERROR_CHECK_EXIST_USER_FAIL:     "检查用户是否存在失败",
	ERROR_EDIT_USER_FAIL:            "修改用户失败",
	ERROR_COUNT_USER_FAIL:           "统计用户失败",
	ERROR_GET_USERS_FAIL:            "获取多个用户失败",
	ERROR_GET_USER_FAIL:             "获取单个用户失败",
	ERROR_GEN_CERT_POSTER_FAIL:      "生成证书海报失败",
	ERROR_AUTH_CHECK_TOKEN_FAIL:     "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT:  "Token已超时",
	ERROR_AUTH_TOKEN:                "Token生成失败",
	ERROR_AUTH:                      "Token错误",
	ERROR_UPLOAD_CREATE_IMAGE_FAIL:  "生成图片失败",
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",
	ERROR_UPLOAD_CREATE_FILE_FAIL:   "生成文件失败",
	ERROR_UPLOAD_SAVE_FILE_FAIL:     "保存文件失败",
	ERROR_UPLOAD_CHECK_FILE_FAIL:    "检查文件失败",
	ERROR_UPLOAD_CHECK_FILE_FORMAT:  "校验文件错误，文件格式或大小有问题",
	ERROR_CERT_APPLY_DISABLED:		 "该证书已关闭申请",
	ERROR_EXIST_APPLY:				 "您已申请过该证书",
	ERROR_CHECK_EXIST_APPLY_FAIL:	 "检测用户是否已申请证书失败",
	ERROR_ADD_APPLY:				 "提交申请失败",
	ERROR_EXIST_APPLY_PAY:			 "您有未处理的支付订单, 请等待退款后再次申领",
	ERROR_EXIST_APPLY_PID:			 "用户已申请过该证书",
	ERROR_GET_USER_CERT_IMAGES:		 "未查询到该用户的电子证书",
	ERROR_GET_USER_CERT_FILES:		 "未查询到该用户的电子证书",
	ERROR_AUTH_CHANGE_PASSWORD_FAIL: "修改密码时发生错误",
	ERROR_AUTH_CHECK_USRNAME_FAIL:   "用户名不存在",
	ERROR_AUTH_CHECK_PASSWORD_FAIL:  "密码不正确",
	ERROR_EXPORT_EMPYT_FILE:         "未查询到可导出的数据",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
