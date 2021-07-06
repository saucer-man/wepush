package main

import (
	"os"
	"path/filepath"

	"github.com/saucerman/wepush/config"
	"github.com/saucerman/wepush/server"

	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func main() {
	config := config.LoadConfig()

	// 设置日志的输出格式和路径
	setLog(config)

	// 运行 http server
	server.Start(&config)
}

func setLog(config config.Config) {
	// 设置日志输出格式
	log.SetFormatter(&nested.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		NoColors:        true,
	})

	// 设置日志输出级别
	logLevel := map[string]log.Level{
		"debug": log.DebugLevel,
		"info":  log.InfoLevel,
		"warn":  log.WarnLevel,
	}
	log.SetLevel(logLevel[config.LogLevel])

	// 把日志输出到文本中
	var logPath = config.LogPath
	os.MkdirAll(logPath, os.ModePerm)
	fileInfo, err := os.Stat(logPath)
	if err == nil && fileInfo.IsDir() {
		logPattern := filepath.Join(logPath, "log.%Y%m%d")
		logCurrent := filepath.Join(logPath, "log.current")
		rl, _ := rotatelogs.New(
			logPattern,
			rotatelogs.WithLinkName(logCurrent),
		)
		log.SetOutput(rl)
	} else {
		log.Warnf("设置日志路径失败，目标路径 %s 不存在或者不是目录", logPath)
	}
}
