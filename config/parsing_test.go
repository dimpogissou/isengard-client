package config

import (
	"errors"
	"testing"
)

// Tests that a valid configuration doesn't return any error at validation
func TestValidS3Config(t *testing.T) {

	var validConnector = []S3ConnectorConfig{S3ConnectorConfig{Name: "somename", Type: "s3", Endpoint: "someEndpoint", KeyPrefix: "prefix", Bucket: "bucket", Region: "region"}}
	var validConfig = YamlConfig{Directory: "./", ConfigName: "something", LogPattern: "something", S3Connectors: validConnector}
	got := validateConfig(validConfig)
	if got != nil {
		t.Errorf("validateConfig(%v) == %v, want %v", validConfig, got, nil)
	}
}

// Tests various invalid configuration cases and asserts over the error returned
func TestInvalidConfig(t *testing.T) {

	var unsupportedTypeConnector = []S3ConnectorConfig{S3ConnectorConfig{Name: "somename", Type: "wrongType", Endpoint: "someEndpoint", KeyPrefix: "prefix", Bucket: "bucket", Region: "region", Levels: []string{"INFO", "WARNING"}}}
	var unsupportedLevelConnector = []S3ConnectorConfig{S3ConnectorConfig{Name: "somename", Type: "s3", Endpoint: "someEndpoint", KeyPrefix: "prefix", Bucket: "bucket", Region: "region", Levels: []string{"INFO", "INVALID"}}}

	cases := []struct {
		in   YamlConfig
		want error
	}{
		{YamlConfig{Directory: ""}, errors.New("Did not find logs directory in YAML configuration")},                                                                                     // Config with empty directory
		{YamlConfig{Directory: "./non_existing_directory_123"}, errors.New("Resolved logs directory ./non_existing_directory_123 does not exist, exiting")},                              // Non existing directory
		{YamlConfig{Directory: "./", ConfigName: ""}, errors.New("YAML configuration missing required 'ConfigName' key, exiting")},                                                       // Missing ConfigName
		{YamlConfig{Directory: "./", ConfigName: "something", LogPattern: ""}, errors.New("YAML configuration missing required 'LogPattern' key, exiting")},                              // Missing LogPattern
		{YamlConfig{Directory: "./", ConfigName: "something", LogPattern: "something", S3Connectors: unsupportedTypeConnector}, errors.New("Invalid connector type: wrongType")},         // Unsupported Connector Type
		{YamlConfig{Directory: "./", ConfigName: "something", LogPattern: "something", S3Connectors: unsupportedLevelConnector}, errors.New("Invalid value for logging level: INVALID")}, // Unsupported Connector Level
	}
	for _, c := range cases {
		got := validateConfig(c.in)
		if got.Error() != c.want.Error() {
			t.Errorf("validateConfig(%v) == %v, want %v", c.in, got, c.want)
		}
	}
}
