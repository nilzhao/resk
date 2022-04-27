package base

import (
	"os"

	"github.com/natefinch/lumberjack"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func init() {
	// 定义日志格式
	formatter := &prefixed.TextFormatter{
		ForceColors:     true, // 控制台高亮
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceFormatting: true,
	}
	// formatter.SetColorScheme(&prefixed.ColorScheme{
	// 	InfoLevelStyle: "yellow",
	// })
	log.SetFormatter(formatter)

	// 日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}

	// 日志文件和滚动配置
	log.SetOutput(&lumberjack.Logger{
		// 日志名称
		Filename: "logs/resk.log",
		// 日志大小限制，单位MB
		MaxSize: 100,
		// 历史日志文件保留天数
		MaxAge: 30,
		// 最大保留历史日志数量
		MaxBackups: 30,
		// 本地时区
		LocalTime: true,
		// 历史日志文件压缩标识
		Compress: false,
	})
}
