package testx

import (
	"resk/infra"
	"resk/infra/base"

	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func init() {
	file := kvs.GetCurrentFilePath("../brun/config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})

	app := infra.New(conf)
	app.Start()
}
