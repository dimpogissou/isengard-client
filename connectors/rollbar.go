package connectors

import (
	"github.com/dimpogissou/isengard-server/config"
	"github.com/hpcloud/tail"
)

type RollbarConnector struct {
	cfg config.Connector
}

func (c RollbarConnector) Open() {
	log.Info("Starting Rollbar connector ...")
}
func (c RollbarConnector) Close() {}

func (c RollbarConnector) Send(line *tail.Line) bool {
	log.Warning("Sending line to Rollbar -->", line)
	return true
}
