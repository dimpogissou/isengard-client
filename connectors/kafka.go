package connectors

import (
	"context"
	"fmt"
	"log"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/segmentio/kafka-go"
)

func SetupKafkaConnection(host string, port string, topic string) *kafka.Writer {

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{fmt.Sprintf("%s:%s", host, port)},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	return writer
}

func CloseKafkaConnection(writer *kafka.Writer) {

	if err := writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}

func (c KafkaConnector) writeKafkaMessages(key string, message string) {

	err := c.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
		})
	if err != nil {
		log.Fatal("Failed to write message:", err)
	}

}

type KafkaConnector struct {
	cfg    KafkaConnectorConfig
	writer *kafka.Writer
}

func (c KafkaConnector) Open() {
}
func (c KafkaConnector) Close() {
}

func (c KafkaConnector) Send(line *tail.Line) bool {
	logger.Info(fmt.Sprintf("Sending line to Kafka --> %v", line.Text))
	uuid, err := uuid.NewV4()
	if err != nil {
		return false
	}
	c.writeKafkaMessages(fmt.Sprintf("%v", uuid), line.Text)
	return true
}
