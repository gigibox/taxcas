package main

import (
	"fmt"
	"github.com/fvbock/endless"
	"gopkg.in/urfave/cli.v1"
	"log"
	"net/http"
	"os"
	"runtime"
	"syscall"
	"taxcas/models"
	"taxcas/pkg/export"
	"taxcas/pkg/gredis"
	"taxcas/pkg/logging"
	"taxcas/pkg/setting"
	"taxcas/pkg/upload"
	"taxcas/routers"
)

// @title TAXCAS Example API
// @version 1.0
// @description Certificate authentication system.

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var configFile string

func runService() {
	setting.Setup(configFile)
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
		log.Fatal("Server err: %v", err)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "taxcas"
	app.Version = "1.0.0"
	app.Usage = "Certificate authentication system."
	app.Flags = []cli.Flag {
			cli.StringFlag{
				Name:  "conf, c",
				Value: "conf/app.ini",
				Usage: "Specify configuration file.",
			},
	}

	app.Action = func(c *cli.Context) error {
		configFile = c.String("conf")
		runService()
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
