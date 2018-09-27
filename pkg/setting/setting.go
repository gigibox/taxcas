package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret string
	PageSize  int
	PrefixUrl string

	RuntimeRootPath string
	UploadAllowMaxSize   int

	ImageAllowExts []string
	ExcelAllowExts []string

	UploadSavePath string
	ExportSavePath string
	QrCodeSavePath string
	FontSavePath string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string
}
var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}
var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}
var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}
var RedisSetting = &Redis{}

type Weixin struct {
	AppID      string
	AppSecret  string
	MchID      string
	ApiKey     string
	Notify_url string
}
var WeixinSetting = &Weixin{}

var cfg *ini.File

func Setup(conf string) {
	var err error
	cfg, err = ini.Load(conf)
	if err != nil {
		log.Fatalf("Fail to parse '%s': %v", conf, err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)
	mapTo("redis", RedisSetting)
	mapTo("weixin", WeixinSetting)

	AppSetting.UploadAllowMaxSize = AppSetting.UploadAllowMaxSize * 1024 * 1024
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
	RedisSetting.IdleTimeout = RedisSetting.IdleTimeout * time.Second
}

func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
