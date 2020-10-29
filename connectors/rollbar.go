package connectors

import (
	"fmt"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

type RollbarConnector struct{ cfg RollbarConnectorConfig }

func (c RollbarConnector) GetName() string {
	return c.cfg.Name
}

func (c RollbarConnector) Close() error {
	return nil
}

func (c RollbarConnector) Send(line *tail.Line) error {
	logger.Info(fmt.Sprintf("Sending line to Rollbar --> %v", line))
	return nil
}
