package connectors

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
	"github.com/segmentio/kafka-go"
)

// reads Kafka topic messages to assert they were correctly published
func readFromTopic(host string, port string, topic string, partition int) []byte {

	// make a new reader that consumes from topic-A
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{fmt.Sprintf("%s:%s", host, port)},
		Partition: partition,
		Topic:     topic,
		MinBytes:  1,    // 10B
		MaxBytes:  10e6, // 10MB
	})

	m, err := r.ReadMessage(context.Background())
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

	if err := r.Close(); err != nil {
		panic(err.Error())
	}

	return m.Value
}

// creates a Kafka topic for testing
func createTopic(host string, port string, topic string, partition int) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", fmt.Sprintf("%s:%s", host, port), topic, partition)
	if err != nil {
		panic("Failed creating Kafka connection")
	}
	defer conn.Close()
	topicConfigs := []kafka.TopicConfig{
		kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	logger.Info(fmt.Sprintf("Creating Kafka topic %s ...", topic))
	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		panic(err.Error())
	}
}

func TestWriteToKafkaTopic(t *testing.T) {

	host := os.Getenv("KAFKA_HOST")
	port := os.Getenv("KAFKA_PORT")
	topic := os.Getenv("KAFKA_TOPIC")
	testMessage := "Test message"
	partition := 0
	createTopic(host, port, topic, partition)

	// Create test S3 bucket
	cfg := KafkaConnectorConfig{
		Name:   "testKafkaConnector",
		Type:   "kafka",
		Levels: []string{"INFO", "DEBUG", "WARN", "ERROR"},
		Host:   host,
		Port:   port,
		Topic:  topic,
	}

	connector := KafkaConnector{cfg: cfg, writer: SetupKafkaConnection(cfg.Host, cfg.Port, cfg.Topic)}
	defer CloseKafkaConnection(connector.writer)

	line := tail.Line{Text: testMessage}

	connector.Send(&line)

	msg := readFromTopic(cfg.Host, cfg.Port, cfg.Topic, partition)

	if string(msg) != testMessage {
		t.Errorf("Kafka message received not matching expected bytes, got %v, want %v", msg, testMessage)
	}

}
