package config

import (
	"errors"
	"fmt"
)

// S3 connector configuration
type S3ConnectorConfig struct {
	Name      string   `yaml:"Name"`
	Endpoint  string   `yaml:"Endpoint"`
	KeyPrefix string   `yaml:"KeyPrefix"`
	Bucket    string   `yaml:"Bucket"`
	Region    string   `yaml:"Region"`
	Type      string   `yaml:"Type"`
	Levels    []string `yaml:"Levels"`
}

func (config S3ConnectorConfig) getName() string {
	return config.Name
}

func (config S3ConnectorConfig) getType() string {
	return config.Type
}

func (config S3ConnectorConfig) getLevels() []string {
	return config.Levels
}

func (config S3ConnectorConfig) validate() error {
	if missingFields(config.Endpoint, config.KeyPrefix, config.Bucket, config.Region) {
		return errors.New(
			fmt.Sprintf("Missing field(s) in S3 connector config '%s': endpoint = %s, keyprefix = %s, bucket = %s, region = %s",
				config.Name, config.Endpoint, config.KeyPrefix, config.Bucket, config.Region))
	}
	return nil
}
