package bootstrap

import (
	"github.com/goecology/webhook/app/pkg/conf"
	"github.com/goecology/webhook/app/pkg/mus"
)

var Arg arg

type arg struct {
	CfgFile string
	Local   bool
}

func Init() {
	mus.Init(Arg.CfgFile)
	conf.Init()
}
