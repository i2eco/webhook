package conf

import (
	"github.com/spf13/viper"
)

var Conf = conf{}

type conf struct {
	WebHook map[string]UrlInfo
}

type UrlInfo struct {
	UrlPath    string
	Token      string
	ExecPath   string
	ExecParams []string
	IsBash     bool
}

func Init() error {
	err := viper.Unmarshal(&Conf)
	return err
}
