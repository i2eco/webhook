package conf

import (
	"github.com/goecology/webhook/app/pkg/mus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

func Init() {
	err := viper.Unmarshal(&Conf)
	if err != nil {
		mus.Logger.Error("unmarshal error", zap.Error(err))
		return
	}
}
