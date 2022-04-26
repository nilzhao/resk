package main

import (
	"fmt"
	"time"

	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func main() {
	file := kvs.GetCurrentFilePath("config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)
	port := conf.GetIntDefault("app.server.port", 18080)
	fmt.Println(port)
	fmt.Println(conf.GetDefault("app.name", "unknown"))
	fmt.Println(conf.GetDurationDefault("app.time", time.Second))
}
