package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strings"
)

var (
	confPath  string
	region    string
	zone      string
	deployEnv string
	host      string
	// Conf config
	Conf *Config
)

type Config struct {
	Token *Token
}

type Token struct {
	Mid      int    `json:"mid"`
	RoomId   string `json:"room_id"`
	Platform string `json:"platform"`
	Accepts  []int  `json:"accepts"`
}

func init() {
	//var (
	//	defHost, _ = os.Hostname()
	//)
	if IsRelease() {
		flag.StringVar(&confPath, "conf", "tcp-example.toml", "default config path")
	} else {
		flag.StringVar(&confPath, "conf", "cmd/tcp/tcp-example.toml", "default config path")
	}
	//flag.StringVar(&region, "region", os.Getenv("REGION"), "avaliable region. or use REGION env variable, value: sh etc.")
	//flag.StringVar(&zone, "zone", os.Getenv("ZONE"), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
	//flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	//flag.StringVar(&host, "host", defHost, "machine hostname. or use default machine hostname.")
}

func Init() (err error) {
	Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func Default() *Config {
	return &Config{
		Token: &Token{
			Platform: "web",
			Accepts:  []int{1000, 1001, 1002},
		},
	}
}

func IsRelease() bool {
	arg1 := strings.ToLower(os.Args[0])
	name := filepath.Base(arg1)

	return strings.Index(name, "__") != 0 && strings.Index(arg1, "go-build") < 0
}
