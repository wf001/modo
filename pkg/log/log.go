package log

import (
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/sirupsen/logrus"
)

var DEFAULT_FORMAT = "%+v"

const (
	FgBlack int = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

func getCaller() (string, string, int) {
	pc, file, line, _ := runtime.Caller(4)
	files := regexp.MustCompile("[/]").Split(file, -1)
	funcName := regexp.MustCompile("[/]").Split(runtime.FuncForPC(pc).Name(), -1)

	return funcName[len(funcName)-1], files[len(files)-1], line
}
func logrusWithField() *logrus.Entry {
	funcName, file, line := getCaller()

	return logrus.WithFields(logrus.Fields{
		"func": funcName,
		"file": file,
		"line": line,
	})
}

func BLUE(body string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", FgCyan, body)
}
func GREEN(body string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", FgGreen, body)
}
func YELLOW(body string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", FgYellow, body)
}

func DebugMessage(message string) {
	debug("", YELLOW(message))
}

func DebugNoLine(format string, value ...interface{}) {
	debugNoLine(format, value...)
}

func debugNoLine(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}
	logrus.Debugf(format, value...)
}
func Debug(format string, value ...interface{}) {
	debug(format, value...)
}

func debug(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}
	logrusWithField().Debugf(format, value...)
}
func DebugValueColored(format string, value interface{}) {
	Debug(
		fmt.Sprintf(format, GREEN(fmt.Sprintf("%s", value))),
	)
}

func Info(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Infof(format, value...)
}

func Warn(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Warnf(format, value...)
}

func Error(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	logrus.Errorf(format, value...)
}

func Panic(format string, value ...interface{}) {
	defaultFormat := DEFAULT_FORMAT
	if format == "" {
		format = defaultFormat
	}

	panic(fmt.Sprintf(format, value...))
}

func init() {
	// TODO: support switching
	logrus.SetOutput(os.Stdout)
	SetLevelWarning()
	SetFormatter()
}

func SetOutputFile(fileName string) {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.SetOutput(f)
	defer f.Close()
}

func SetFormatter() {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
}

func SetLevelError() {
	logrus.SetLevel(logrus.ErrorLevel)
}
func SetLevelWarning() {
	logrus.SetLevel(logrus.WarnLevel)
}
func SetLevelInfo() {
	logrus.SetLevel(logrus.InfoLevel)
}
func SetLevelDebug() {
	logrus.SetLevel(logrus.DebugLevel)
}
