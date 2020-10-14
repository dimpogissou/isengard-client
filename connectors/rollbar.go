package connectors

import (
	"fmt"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

type RollbarConnector struct{ cfg RollbarConnectorConfig }

func (c RollbarConnector) Open() {
	logger.Info("Starting Rollbar connector ...")
}
func (c RollbarConnector) Close() {}

func (c RollbarConnector) Send(line *tail.Line) bool {
	logger.Info(fmt.Sprintf("Sending line to Rollbar --> %v", line))
	return true
}
