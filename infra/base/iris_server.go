package base

import (
	"resk/infra"
	"time"

	"github.com/kataras/iris/v12"
	irisLogger "github.com/kataras/iris/v12/middleware/logger"
	irisRecover "github.com/kataras/iris/v12/middleware/recover"
	"github.com/sirupsen/logrus"
)

var irisApplication *iris.Application

type IrisServerStarter struct {
	infra.BaseStarter
}

func Iris() *iris.Application {
	return irisApplication
}

func (s *IrisServerStarter) Init(ctx infra.StarterContext) {
	// 初始化 iris
	irisApplication = initIris()
	// 日志组件扩展
	logger := irisApplication.Logger()
	logger.Install(logrus.StandardLogger())
}

func (s *IrisServerStarter) Start(ctx infra.StarterContext) {
	routes := Iris().GetRoutes()
	for _, route := range routes {
		logrus.Info(route.Method + route.Path)
	}
	port := ctx.Props().GetDefault("app.server.port", "18080")
	err := Iris().Listen(":" + port)
	if err != nil {
		panic(err)
	}
}

func (s *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	app.UseRouter(irisRecover.New())
	conf := irisLogger.Config{
		Status:     true,
		IP:         true,
		Method:     true,
		Path:       true,
		Query:      true,
		TraceRoute: true,
		LogFunc: func(endTime time.Time, latency time.Duration, status, ip, method, path string, message, headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s |", endTime.Format("2006-01-02 15:04:05"), latency.String(), status, ip, method, path, message, headerMessage)
		},
	}
	app.Use(irisLogger.New(conf))
	return app
}
