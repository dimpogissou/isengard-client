package logger

import (
	"fmt"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("mainLogger")

// Initialise logging configuration
func InitLogger() {

	var backend = logging.NewLogBackend(os.Stderr, "", 0)
	var format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
	var _ = logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backend)

}

func Debug(msg string) {
	log.Debug(msg)
}

func Info(msg string) {
	log.Info(msg)
}

func Warn(alertCode string, msg string) {
	log.Warning(fmt.Sprintf("[%s] %s", alertCode, msg))
}

func CheckWarnAndLog(err error, alertCode string, msg string) {
	if err != nil {
		log.Warning(fmt.Sprintf("[%s] %s - Error --> %s", alertCode, msg, err.Error()))
	}
}

func Error(alertCode string, msg string) {
	log.Error(fmt.Sprintf("[%s] %s", alertCode, msg))
}

func CheckErrAndLog(err error, alertCode string, msg string) {
	if err != nil {
		log.Error(fmt.Sprintf("[%s] %s - Error --> %s", alertCode, msg, err.Error()))
	}
}

func CheckErrAndPanic(err error, alertCode string, msg string) {
	if err != nil {
		log.Error(fmt.Sprintf("[%s] %s - Error --> %s", alertCode, msg, err.Error()))
		panic(err.Error)
	}
}
