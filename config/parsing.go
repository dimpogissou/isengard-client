package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/dimpogissou/isengard-server/logger"
	"gopkg.in/yaml.v2"
)

var supportedConnectors = []string{"s3", "rollbar", "kafka"}
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
	validate() error
}

func getConnectorsConfigs(cfg YamlConfig) []ConnectorConfig {

	connectorsConfigs := []ConnectorConfig{}
	for _, connCfg := range cfg.S3Connectors {
		connectorsConfigs = append(connectorsConfigs, connCfg)
	}
	for _, connCfg := range cfg.KafkaConnectors {
		connectorsConfigs = append(connectorsConfigs, connCfg)
	}
	for _, connCfg := range cfg.RollbarConnectors {
		connectorsConfigs = append(connectorsConfigs, connCfg)
	}

	return connectorsConfigs

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

	connectorsConfigs := getConnectorsConfigs(cfg)

	for _, connCfg := range connectorsConfigs {
		// Assert connectors have valid common fields values
		err := validateConnectorsCommonFields(connCfg)
		if err != nil {
			return err
		}
		// Assert connectors have valid S3 connector fields
		err = connCfg.validate()
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
	if !stringInSlice(connector.getType(), supportedConnectors) {
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
