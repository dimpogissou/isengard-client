package connectors

import (
	"github.com/hpcloud/tail"
)

type ConnectorInterface interface {
	Open()
	Send(line *tail.Line) bool
	Close()
}

// Create all connectors
func GenerateConnectors(cfg YamlConfig) []ConnectorInterface {

	conns := []ConnectorInterface{}

	for _, connCfg := range cfg.S3Connectors {
		conns = append(conns, S3Connector{cfg: connCfg, client: SetupS3Client(connCfg)})
	}

	for _, connCfg := range cfg.RollbarConnectors {
		conns = append(conns, RollbarConnector{cfg: connCfg})
	}

	return conns
}
