package connectors

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/dimpogissou/isengard-server/logger"
	"gopkg.in/yaml.v2"
)

var SupportedConnectors = []string{"s3", "rollbar", "kafka"}
var supportedLevels = []string{"DEBUG", "INFO", "WARNING", "WARN", "ERROR"}

// YAML configuration structs
type YamlConfig struct {
	ConfigName        string                   `yaml:"ConfigName"`
	Directory         string                   `yaml:"Directory"`
	LogPattern        string                   `yaml:"LogPattern"`
	Definitions       []PatternConfig          `yaml:"Definitions"`
	S3Connectors      []S3ConnectorConfig      `yaml:"S3Connectors"`
	RollbarConnectors []RollbarConnectorConfig `yaml:"RollbarConnectors"`
	KafkaConnectors   []KafkaConnectorConfig   `yaml:"KafkaConnectors"`
}

type PatternConfig struct {
	Name    string `yaml:"Name"`
	Pattern string `yaml:"Pattern"`
}

type ConnectorConfig interface {
	getName() string
	getType() string
	getLevels() []string
}

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

func validateS3ConnectorsFields(connector S3ConnectorConfig) error {
	if missingFields(connector.Endpoint, connector.KeyPrefix, connector.Bucket, connector.Region) {
		return errors.New(
			fmt.Sprintf("Missing field(s) in S3 connector config '%s': endpoint = %s, keyprefix = %s, bucket = %s, region = %s",
				connector.Name, connector.Endpoint, connector.KeyPrefix, connector.Bucket, connector.Region))
	}
	return nil
}

// Rollbar connector configuration
type RollbarConnectorConfig struct {
	Name   string   `yaml:"Name"`
	Type   string   `yaml:"Type"`
	Url    string   `yaml:"Url"`
	Levels []string `yaml:"Levels"`
}

func (config RollbarConnectorConfig) getName() string {
	return config.Name
}

func (config RollbarConnectorConfig) getType() string {
	return config.Type
}

func (config RollbarConnectorConfig) getLevels() []string {
	return config.Levels
}

func validateRollbarConnectorsFields(connector RollbarConnectorConfig) error {
	return nil
}

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

func validateKafkaConnectorsFields(connector KafkaConnectorConfig) error {
	if missingFields(connector.Host, connector.Port, connector.Topic) {
		return errors.New(
			fmt.Sprintf("Missing field(s) in S3 connector config '%s': host = %s, port = %s, topic = %s",
				connector.Name, connector.Host, connector.Port, connector.Topic))
	}
	return nil
}

// Runs complete configuration validation steps and returns eventual errors
func validateConfig(cfg YamlConfig) error {

	// If Directory is nil, try to retrieve from env var, if not then error
	if cfg.Directory == "" {
		return errors.New("Did not find logs directory in YAML configuration")
	}

	// If Directory resolved, check if directory exists, if not then error
	if _, err := os.Stat(cfg.Directory); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("Resolved logs directory %s does not exist, exiting", cfg.Directory))
	}

	// If Name or LogPattern missing, error
	if cfg.ConfigName == "" {
		return errors.New("YAML configuration missing required 'ConfigName' key, exiting")
	} else if cfg.LogPattern == "" {
		return errors.New("YAML configuration missing required 'LogPattern' key, exiting")
	}

	for _, connector := range cfg.S3Connectors {
		// Assert connectors have valid common fields values
		err := validateConnectorsCommonFields(connector)
		if err != nil {
			return err
		}
		// Assert connectors have valid S3 connector fields
		err = validateS3ConnectorsFields(connector)
		if err != nil {
			return err
		}
	}

	for _, connector := range cfg.RollbarConnectors {
		// Assert connectors have valid common fields values
		err := validateConnectorsCommonFields(connector)
		if err != nil {
			return err
		}
		// Assert connectors have valid Rollbar connector fields
		err = validateRollbarConnectorsFields(connector)
		if err != nil {
			return err
		}
	}

	for _, connector := range cfg.KafkaConnectors {
		// Assert connectors have valid common fields values
		err := validateConnectorsCommonFields(connector)
		if err != nil {
			return err
		}
		// Assert connectors have valid Kafka connector fields
		err = validateKafkaConnectorsFields(connector)
		if err != nil {
			return err
		}
	}

	return nil
}

// Util function returning true if any of the provided strings is empty, false otherwise
func missingFields(fields ...string) bool {
	for _, field := range fields {
		if field == "" {
			return true
		}
	}
	return false
}

// Validates common fields for all monitors: Name, Type, Levels
func validateConnectorsCommonFields(connector ConnectorConfig) error {

	if missingFields(connector.getName(), connector.getType()) {
		return errors.New(fmt.Sprintf("Missing field in connector config: %v", connector))
	}
	if !stringInSlice(connector.getType(), SupportedConnectors) {
		return errors.New(fmt.Sprintf("Invalid connector type: %s", connector.getType()))
	}
	for _, level := range connector.getLevels() {
		if !stringInSlice(level, supportedLevels) {
			return errors.New(fmt.Sprintf("Invalid value for logging level: %s", level))
		}
	}
	return nil
}

// Reads and parses YAML configuration
func readConfig(path string) YamlConfig {

	conf := YamlConfig{}
	data, readErr := ioutil.ReadFile(path)

	// If error at file read, log and stop execution
	if readErr != nil {
		logger.CheckErrAndPanic(readErr, "FailedReadingConfigFile", fmt.Sprintf("Could not read YAML configuration at %s", path))
	}

	// If error at file parsing, log and stop
	parseErr := yaml.Unmarshal(data, &conf)
	if parseErr != nil {
		logger.CheckErrAndPanic(parseErr, "FailedReadingConfigFile", "Error occurred while parsing YAML configuration file")
	}

	return conf
}

// Builds Regex specified in configuration
func BuildRegex(cfg YamlConfig) *regexp.Regexp {

	// Create subPatterns slice from cfg.Definitions
	subPatterns := make([]interface{}, len(cfg.Definitions))
	for i, def := range cfg.Definitions {
		subPatterns[i] = def.Pattern
	}

	// Interpolate subpatterns in main pattern, compile regex
	pattern := fmt.Sprintf(cfg.LogPattern, subPatterns...)
	regex := regexp.MustCompile(pattern)

	return regex

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Core function parsing, validating and returning YAML config
func ValidateAndLoadConfig(path *string) YamlConfig {

	cfg := readConfig(*path)
	err := validateConfig(cfg)
	logger.CheckErrAndPanic(err, "FailedValidatingConfigFile", fmt.Sprintf("Configuration file validation failed for %s", *path))

	return cfg

}
