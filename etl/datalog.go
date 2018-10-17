package etl

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// 获取命令行日志
func GetLog() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = &logrus.TextFormatter{TimestampFormat: "20060102150405"}
	logger.Out = os.Stdout
	logger.Level = logrus.DebugLevel
	return logger
}

// 获取文件日志
func GetFileLog(appinfo AppInfo) (*logrus.Logger, error) {
	logger := GetLog()
	logpath := filepath.Join(appinfo.Home, "log", appinfo.Name+".log")
	logfile, err := os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, errors.New("加载日志文件出错" + logpath)
	}
	logger.Out = logfile
	return logger, nil
}
