package connectors

import (
	"context"
	"fmt"

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

func CloseKafkaConnection(writer *kafka.Writer) error {

	if err := writer.Close(); err != nil {
		logger.Error("KafkaClosePublisherError:", err.Error())
		return err
	}
	logger.Info("Closed Kafka publisher ...")
	return nil
}

func (c KafkaConnector) writeKafkaMessages(key string, message string) error {

	err := c.writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(key),
			Value: []byte(message),
		})
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Successfully published message to Kafka at key [%s] -> %s", key, message))
	return nil
}

type KafkaConnector struct {
	cfg    KafkaConnectorConfig
	writer *kafka.Writer
}

func (c KafkaConnector) Close() {
	CloseKafkaConnection(c.writer)
}

func (c KafkaConnector) Send(line *tail.Line) error {
	logger.Debug(fmt.Sprintf("Sending line to Kafka --> %v", line.Text))
	uuid, err := uuid.NewV4()
	if err != nil {
		logger.Error("CreateUuidError", err.Error())
		return err
	}
	// TODO -> Optimize string write since this operation is repeated for each log line
	err = c.writeKafkaMessages(fmt.Sprintf("%v", uuid), line.Text)
	if err != nil {
		logger.Error("KafkaPublishMessageError", err.Error())
		return err
	}
	return nil
}
