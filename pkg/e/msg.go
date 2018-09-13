package e

var MsgFlags = map[int]string{
	SUCCESS:                         "ok",
	ERROR:                           "fail",
	INVALID_PARAMS:                  "请求参数错误",
	ERROR_EXIST_CERT:                "已存在该证书名称",
	ERROR_EXIST_CERT_FAIL:           "获取已存在证书失败",
	ERROR_NOT_EXIST_CERT:            "该证书不存在",
	ERROR_GET_CERTS_FAIL:            "获取所有证书失败",
	ERROR_COUNT_CERT_FAIL:           "统计证书失败",
	ERROR_ADD_CERT_FAIL:             "新增证书失败",
	ERROR_EDIT_CERT_FAIL:            "修改证书失败",
	ERROR_DELETE_CERT_FAIL:          "删除证书失败",
	ERROR_EXPORT_CERT_FAIL:          "导出证书失败",
	ERROR_IMPORT_CERT_FAIL:          "导入证书失败",
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
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
