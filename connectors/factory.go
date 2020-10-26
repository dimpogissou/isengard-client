package connectors

import (
	"github.com/hpcloud/tail"
)

type ConnectorInterface interface {
	Send(line *tail.Line) error
	Close() error
}

// Create all connectors
func CreateConnectors(cfg YamlConfig) []ConnectorInterface {

	conns := []ConnectorInterface{}

	for _, connCfg := range cfg.S3Connectors {
		session, client := SetupS3Client(connCfg)
		conns = append(conns, S3Connector{cfg: connCfg, session: session, client: client})
	}

	for _, connCfg := range cfg.RollbarConnectors {
		conns = append(conns, RollbarConnector{cfg: connCfg})
	}

	for _, connCfg := range cfg.KafkaConnectors {
		writer := SetupKafkaConnection(connCfg.Host, connCfg.Port, connCfg.Topic)
		conns = append(conns, KafkaConnector{cfg: connCfg, writer: writer})
	}

	return conns
}
