package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

var Conf config

type log struct {
	LogDir         string `toml:"log_dir"`
	LogFileName    string `toml:"log_file_name"`
	LogLevel       string `toml:"log_level"`
	MaxSize        int    `toml:"max_size"`
	MaxAge         int    `toml:"max_age"`
	MaxBackups     int    `toml:"max_backups"`
}

type redisKey struct {
	Publisher string `toml:"publisher"`
	PayList   string `toml:"pay_list"`
	VIPList   string `toml:"vip_list"`
	Member    string `toml:"member"`
	VIP       string `toml:"vip"`
}

type redis struct {
	Host         string        `toml:"host"`
	Password     string        `toml:"password"`
	Key        redisKey
}

type mysql struct {
	User     string `toml:"user"`
	Password string `toml:"password"`
	Type     string `toml:"type"`
	Path     string `toml:"path"`
	DBName   string `toml:"dbname"`
}

type config struct {
	Log     log
	Redis   redis
	Mysql   mysql
}

func init() {
	if _, err := toml.DecodeFile("conf/config.toml", &Conf); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
