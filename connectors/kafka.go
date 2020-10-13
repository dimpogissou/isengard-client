package connectors

import (
	"github.com/dimpogissou/isengard-server/config"
	"github.com/hpcloud/tail"
)

type KafkaConnector struct {
	cfg config.Connector
}

func (c KafkaConnector) Open()  {}
func (c KafkaConnector) Close() {}

func (c KafkaConnector) Send(line *tail.Line) bool {
	log.Warning("Sending line to Kafka -->", line)
	return true
}
