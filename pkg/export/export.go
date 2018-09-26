package export

import (
	"taxcas/pkg/setting"
	"taxcas/pkg/util"
)

func Setup() {
	util.CheckDir(GetExportFullPath())
	util.CheckDir(GetExportFullPath() + "image/")
	util.CheckDir(GetExportFullPath() + "pdf/")
	util.CheckDir(GetExportFullPath() + "excel/")
}

func GetExportPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetRuntimePath() string {
	return setting.AppSetting.RuntimeRootPath
}

func GetExportFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExportPath()
}

func GetExportImagePath(certid string) string {
	return GetExportPath() + "image/" + "cert" + certid + "/"
}

func GetExportPDFPath(certid string) string {
	return GetExportPath() + "pdf/" + "cert" + certid + "/"
}

func GetExportExcelPath() string {
	return GetExportPath() + "excel/"
}