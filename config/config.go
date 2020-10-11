package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("main")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var supportedConnectors = []string{"s3", "rollbar"}
var supportedLevels = []string{"DEBUG", "INFO", "WARNING", "WARN", "ERROR"}

type PatternConfig struct {
	Name    string `yaml:"Name"`
	Pattern string `yaml:"Pattern"`
}

type Connector struct {
	Name   string   `yaml:"Name"`
	Type   string   `yaml:"Type"`
	Levels []string `yaml:"Levels"`
}

type YamlConfig struct {
	ConfigName  string          `yaml:"ConfigName"`
	Directory   string          `yaml:"Directory"`
	LogPattern  string          `yaml:"LogPattern"`
	Definitions []PatternConfig `yaml:"Definitions"`
	Connectors  []Connector     `yaml:"Connectors"`
}

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

	// Assert Connectors have valid values
	for _, connector := range cfg.Connectors {
		if !stringInSlice(connector.Type, supportedConnectors) {
			return errors.New(fmt.Sprintf("Invalid connector type: %s", connector.Type))
		}
		for _, level := range connector.Levels {
			if !stringInSlice(level, supportedLevels) {
				return errors.New(fmt.Sprintf("Invalid value for logging level: %s", level))
			}
		}
	}

	return nil
}

func readConfig(path string) YamlConfig {

	conf := YamlConfig{}
	data, readErr := ioutil.ReadFile(path)

	// If error at file read, log and stop execution
	if readErr != nil {
		log.Fatalf("Could not read YAML configuration at %s due to: %s", path, readErr)
	}

	// If error at file parsing, log and stop
	parseErr := yaml.Unmarshal(data, &conf)
	if parseErr != nil {
		log.Fatalf("Error occurred while parsing YAML configuration file: %v", parseErr)
	}

	return conf
}

func BuildRegex(cfg YamlConfig) *regexp.Regexp {

	// Create subPatterns slice from cfg.Definitions
	subPatterns := make([]interface{}, len(cfg.Definitions))
	for i, def := range cfg.Definitions {
		subPatterns[i] = def.Pattern
	}

	// Interpolate subpatterns in main pattern, compile regex
	pattern := fmt.Sprintf(cfg.LogPattern, subPatterns...)
	regex := regexp.MustCompile(pattern)

	fmt.Printf("Successfully concatenated log line regular expression --> %s\n", pattern)

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

func ValidateAndLoadConfig() YamlConfig {

	cfgPath := os.Getenv("ISENGARD_CONFIG_FILE")
	if cfgPath == "" {
		log.Fatalf("ISENGARD_CONFIG_FILE environment variable not set")
	}

	cfg := readConfig(cfgPath)
	err := validateConfig(cfg)

	if err != nil {
		log.Fatalf("Configuration validation failed due to: %v", err)
	}

	return cfg

}
