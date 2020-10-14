package connectors

import (
	"context"
	"fmt"
	"testing"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/segmentio/kafka-go"
)

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

}
