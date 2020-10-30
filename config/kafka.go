package config

import (
	"errors"
	"fmt"
)

// Kafka connector configuration
type KafkaConnectorConfig struct {
	Name   string   `yaml:"Name"`
	Type   string   `yaml:"Type"`
	Host   string   `yaml:"Host"`
	Port   string   `yaml:"Port"`
	Topic  string   `yaml:"Topic"`
	Levels []string `yaml:"Levels"`
}

func (config KafkaConnectorConfig) getName() string {
	return config.Name
}

func (config KafkaConnectorConfig) getType() string {
	return config.Type
}

func (config KafkaConnectorConfig) getLevels() []string {
	return config.Levels
}

func (config KafkaConnectorConfig) validate() error {
	if missingFields(config.Host, config.Port, config.Topic) {
		return errors.New(
			fmt.Sprintf("Missing field(s) in S3 connector config '%s': host = %s, port = %s, topic = %s",
				config.Name, config.Host, config.Port, config.Topic))
	}
	return nil
}
