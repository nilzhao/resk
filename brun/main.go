package main

import (
	"resk/infra"
	"resk/infra/base"

	_ "resk"

	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog()
	app := infra.New(conf)
	app.Start()
}
