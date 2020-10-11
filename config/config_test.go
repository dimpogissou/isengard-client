package config

import (
	"errors"
	"testing"
)

func TestValidConfig(t *testing.T) {

	var validConnector = []Connector{Connector{Name: "somename", Type: "s3"}}
	var validConfig = YamlConfig{Directory: "./", ConfigName: "something", LogPattern: "something", Connectors: validConnector}
	got := validateConfig(validConfig)
	if got != nil {
		t.Errorf("validateConfig(%v) == %v, want %v", validConfig, got, nil)
	}
}

func TestInvalidConfig(t *testing.T) {

	var unsupportedTypeConnector = []Connector{Connector{Name: "somename", Type: "wrongType"}}
	var unsupportedLevelConnector = []Connector{Connector{Name: "somename", Type: "s3", Levels: []string{"INFO", "INVALID"}}}

	cases := []struct {
		in   YamlConfig
		want error
	}{
		{YamlConfig{Directory: ""}, errors.New("Did not find logs directory in YAML configuration")},                                                                                   // Config with empty directory
		{YamlConfig{Directory: "./non_existing_directory_123"}, errors.New("Resolved logs directory ./non_existing_directory_123 does not exist, exiting")},                            // Non existing directory
		{YamlConfig{Directory: "./", ConfigName: ""}, errors.New("YAML configuration missing required 'ConfigName' key, exiting")},                                                     // Missing ConfigName
		{YamlConfig{Directory: "./", ConfigName: "something", LogPattern: ""}, errors.New("YAML configuration missing required 'LogPattern' key, exiting")},                            // Missing LogPattern
		{YamlConfig{Directory: "./", ConfigName: "something", LogPattern: "something", Connectors: unsupportedTypeConnector}, errors.New("Invalid connector type: wrongType")},         // Unsupported Connector Type
		{YamlConfig{Directory: "./", ConfigName: "something", LogPattern: "something", Connectors: unsupportedLevelConnector}, errors.New("Invalid value for logging level: INVALID")}, // Unsupported Connector Level
	}
	for _, c := range cases {
		got := validateConfig(c.in)
		if got.Error() != c.want.Error() {
			t.Errorf("validateConfig(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}
