package base

import (
	"resk/infra"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
)

var database *dbx.Database

func DbxDatabase() *dbx.Database {
	return database
}

type DbxDatabaseStarter struct {
	infra.BaseStarter
}

func (s *DbxDatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	settings := dbx.Settings{}
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}
	logrus.Infof("%+v\n", settings)
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())
	dbx, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	logrus.Info(dbx.Ping())
	database = dbx
}
