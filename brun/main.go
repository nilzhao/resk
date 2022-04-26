package main

import (
	"resk/infra"

	_ "resk"

	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	app := infra.New(conf)
	app.Start()

	c := make(chan int, 1)
	<-c
}
