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

func Error(alertCode string, msg string) {
	log.Error(fmt.Sprintf("[%s] %s", alertCode, msg))
}
