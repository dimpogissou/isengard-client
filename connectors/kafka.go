package connectors

import (
	"fmt"
	"log"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
	"github.com/segmentio/kafka-go"
)

func setupKafkaConnection(host string, port string, topic string) *kafka.Writer {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{fmt.Sprintf("%s:%s", host, port)},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return writer
}

func closeKafkaConnection(writer *kafka.Writer) {

	if err := writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}

type KafkaConnector struct {
	cfg    KafkaConnectorConfig
	writer *kafka.Writer
}

func (c KafkaConnector) Open() {
	c.writer = setupKafkaConnection(c.cfg.Host, c.cfg.Port, c.cfg.Topic)
}
func (c KafkaConnector) Close() {
	c.writer.Close()
}

func (c KafkaConnector) Send(line *tail.Line) bool {
	logger.Info(fmt.Sprintf("Sending line to Kafka --> %v", line))
	return true
}
