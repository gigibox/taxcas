package upload

import (
	"mime/multipart"
	"strings"

	"taxcas/pkg/file"
	"taxcas/pkg/logging"
	"taxcas/pkg/setting"
	"taxcas/pkg/util"
)

func Setup() {
	util.CheckDir(GetImageFullPath())
	util.CheckDir(GetExcelFullPath())
}

func GetImagePath() string {
	return setting.AppSetting.UploadSavePath + "images/"
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckFileSize(f multipart.File) bool {
	size, err := file.GetSize(f)
	if err != nil {
		logging.Warn(err)
		return false
	}

	return size <= setting.AppSetting.UploadAllowMaxSize
}

func GetExcelPath() string {
	return setting.AppSetting.UploadSavePath + "excels/"
}

func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}

func CheckExcelExt(fileName string) bool {
	ext := file.GetExt(fileName)
	for _, allowExt := range setting.AppSetting.ExcelAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}