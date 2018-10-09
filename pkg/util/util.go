package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"taxcas/pkg/file"
	"taxcas/pkg/logging"
	"taxcas/pkg/setting"
	"time"
)

func GetAppFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/api/" + name
}

func RandomStrings(n int) string {
	const (
		letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func GetRandomFileName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(RandomStrings(16), ext)

	return fileName + ext
}

func CheckDir(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	if err = file.IsNotExistMkDir(dir + "/" + src); err != nil {
		if err = os.Mkdir(dir + "/" + src, os.ModePerm); err != nil {
			return fmt.Errorf("mkdir failed! %v", err)
		}
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}

func EncodeMD5(value string) string {
	m := md5.New()
	m.Write([]byte(value))

	return hex.EncodeToString(m.Sum(nil))
}

func GetFontsList() []string {
	var fonts []string

	list, err := ioutil.ReadDir(setting.AppSetting.RuntimeRootPath + setting.AppSetting.FontSavePath)
	if err != nil {
		logging.Error(err)
		return fonts
	}

	for _, v := range list {
		fonts = append(fonts, v.Name())
	}

	return fonts
}

func CompressStr(str string) string {
	if str == "" {
		return ""
	}
	//匹配一个或多个空白符的正则表达式
	reg := regexp.MustCompile("\\s+")
	return reg.ReplaceAllString(str, "")
}
