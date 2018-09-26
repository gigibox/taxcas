package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"net/http"
	"runtime"
	"syscall"
	"taxcas/pkg/export"
	"taxcas/pkg/gredis"
	"taxcas/pkg/logging"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"

	"taxcas/models"
	"taxcas/routers"
)

// @title TAXCAS Example API
// @version 1.0
// @description Certificate authentication system.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	setting.Setup()
	models.Setup()
	logging.Setup()
	gredis.Setup()
	upload.Setup()
	export.Setup()

	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	if runtime.GOOS == "windows" {
		server := &http.Server{
			Addr:           endPoint,
			Handler:        routersInit,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			MaxHeaderBytes: maxHeaderBytes,
		}

		server.ListenAndServe()
		return
	}

	endless.DefaultReadTimeOut = readTimeout
	endless.DefaultWriteTimeOut = writeTimeout
	endless.DefaultMaxHeaderBytes = maxHeaderBytes
	server := endless.NewServer(endPoint, routersInit)
	server.BeforeBegin = func(add string) {
		logging.Info("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		logging.Error("Server err: %v", err)
	}

}
