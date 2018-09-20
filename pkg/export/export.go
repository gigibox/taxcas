package export

import "taxcas/pkg/setting"

func GetExportFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExportPath() + name
}

func GetExportPath() string {
	return setting.AppSetting.ExportSavePath
}

func GetExportFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExportPath()
}
