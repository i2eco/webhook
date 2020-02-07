package main

import (
	"github.com/goecology/muses"
	"github.com/goecology/muses/pkg/cmd"
	musgin "github.com/goecology/muses/pkg/server/gin"
	"github.com/goecology/muses/pkg/server/stat"
	"github.com/goecology/webhook/app/pkg/conf"
	"github.com/goecology/webhook/app/pkg/mus"
	"github.com/goecology/webhook/app/router"
)

func main() {
	app := muses.Container(
		cmd.Register,
		stat.Register,
		musgin.Register,
	)
	app.SetRouter(router.InitRouter)
	app.PreRun(mus.Init, conf.Init)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
