package config

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

func (config RollbarConnectorConfig) validate() error {
	return nil
}
