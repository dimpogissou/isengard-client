package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/hpcloud/tail"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v2"
)

var log = logging.MustGetLogger("main")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var supportedConnectors = []string{"s3", "rollbar"}
var supportedLevels = []string{"DEBUG", "INFO", "WARNING", "WARN", "ERROR"}

func parseLine(l *tail.Line, re *regexp.Regexp) map[string]string {
	match := re.FindStringSubmatch(l.Text)
	if match == nil {
		log.Warning("Wrongly formatted line, returning empty map")
		return make(map[string]string)
	} else {
		paramsMap := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i > 0 && i <= len(match) {
				paramsMap[name] = match[i]
			}
		}
		return paramsMap
	}
}

func tailFile(path string, ch chan *tail.Line) {

	log.Info(fmt.Sprintf("Start tailing file %s", path))

	t, err := tail.TailFile(path, tail.Config{Follow: true, MustExist: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}, ReOpen: true, Poll: true})
	defer t.Stop()

	if err != nil {
		log.Error(fmt.Sprintf("Could not tail file [%s] due to -> %s", path, err))
	}

	for line := range t.Lines {
		ch <- line
	}
}

func tailDirectory(dir string, regex *regexp.Regexp) int {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error(fmt.Sprintf("Could not get files from directory %s due to -> %s", dir, err))
		return 1
	}

	ch := make(chan *tail.Line)
	defer close(ch)

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dir, file.Name())
		go tailFile(filePath, ch)
	}

	for logLine := range ch {

		// Parse lines
		logMap := parseLine(logLine, regex)

		// Filtering out empty maps from invalid lines
		if len(logMap) > 0 {
			log.Info("Received new log line ->", logMap)
		}
	}

	return 0
}

type PatternConfig struct {
	Name    string `yaml:"Name"`
	Pattern string `yaml:"Pattern"`
}

type Connector struct {
	Name    string   `yaml:"Name"`
	Type    string   `yaml:"Type"`
	Address string   `yaml:"Address"`
	Levels  []string `yaml:"Levels"`
}

type YamlConfig struct {
	ConfigName  string          `yaml:"ConfigName"`
	Directory   string          `yaml:"Directory"`
	LogPattern  string          `yaml:"LogPattern"`
	Definitions []PatternConfig `yaml:"Definitions"`
	Connectors  []Connector     `yaml:"Connectors"`
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func validateConfig(cfg YamlConfig) {

	// If Directory is nil, try to retrieve from env var, if not then error
	if cfg.Directory == "" {
		dir := os.Getenv("ISENGARD_LOGS_DIRECTORY")
		if dir == "" {
			log.Fatalf("Did not find logs directory in YAML configuration nor in environment variables, exiting")
		} else {
			cfg.Directory = dir
		}
	}

	// If Directory resolved, check if directory exists, if not then error
	if _, err := os.Stat(cfg.Directory); os.IsNotExist(err) {
		log.Fatalf("Resolved logs directory %s does not exist, exiting", cfg.Directory)
	}

	// If Name or LogPattern missing, error
	if cfg.ConfigName == "" {
		log.Fatalf("YAML configuration missing required 'ConfigName' key, exiting")
	} else if cfg.LogPattern == "" {
		log.Fatalf("YAML configuration missing required 'LogPattern' key, exiting")
	}

	// Assert Connectors have valid values
	for _, connector := range cfg.Connectors {
		if !stringInSlice(connector.Type, supportedConnectors) {
			log.Fatalf("Invalid connector type: %s", connector.Type)
		}
		for _, level := range connector.Levels {
			if !stringInSlice(level, supportedLevels) {
				log.Fatalf("Invalid value for logging level: %s", level)
			}
		}
	}

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

func buildRegex(cfg YamlConfig) *regexp.Regexp {

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

func main() {

	cfgPath := os.Getenv("ISENGARD_CONFIG_FILE")

	cfg := readConfig(cfgPath)
	validateConfig(cfg)

	regex := buildRegex(cfg)

	tailDirectory(cfg.Directory, regex)

}
