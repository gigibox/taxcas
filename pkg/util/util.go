package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"taxcas/pkg/file"
	"taxcas/pkg/setting"
)

func GetAppFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/api/" + name
}

func RandomStrings(n int) string {
	letterBytes := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
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
