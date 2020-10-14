package connectors

import (
	"fmt"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

type KafkaConnector struct{}

func (c KafkaConnector) Open()  {}
func (c KafkaConnector) Close() {}

func (c KafkaConnector) Send(line *tail.Line) bool {
	logger.Info(fmt.Sprintf("Sending line to Kafka --> %v", line))
	return true
}
