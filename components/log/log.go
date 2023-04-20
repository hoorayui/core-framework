package log

// log
import (
	"fmt"
	"io"
	"os"

	"github.com/hoorayui/core-framework/util"
	"github.com/hoorayui/core-framework/util/flag"
	"github.com/sirupsen/logrus"
)

// Logger 日志封装
type Logger struct {
	Logger *logrus.Logger
}

// UTCFormatter 配置logrus时区
type UTCFormatter struct {
	logrus.Formatter
}

// Format 设置时区
func (u UTCFormatter) Format(e *logrus.Entry) ([]byte, error) {
	// 东八区
	e.Time = util.Now()
	return u.Formatter.Format(e)
}

// GetMultiWriter 设置日志写入流
func GetMultiWriter() io.Writer {
	// TODO 替换为应用名称
	path := fmt.Sprintf("%s/app.log", flag.LogDir)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0777)
	if nil != err {
		logrus.Fatal("打开日志文件[%s]失败", path, flag.LogDir)
		return nil
	}
	return io.MultiWriter(os.Stdout, f)
}

func Debug(args ...interface{}) {
	instance.logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	instance.logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	instance.logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	instance.logger.Infof(format, args...)
}

func Error(args ...interface{}) {
	instance.logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	instance.logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	instance.logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	instance.logger.Fatalf(format, args...)
}

func Panic(args ...interface{}) {
	instance.logger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	instance.logger.Panicf(format, args...)
}

func Trace(args ...interface{}) {
	instance.logger.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	instance.logger.Tracef(format, args...)
}

// ErrorHook ...
type ErrorHook struct{}

func (h *ErrorHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ErrorHook) Fire(entry *logrus.Entry) error {
	// entry.Logger.
	return nil
}
