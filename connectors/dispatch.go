package connectors

import (
	"github.com/dimpogissou/isengard-server/config"
	"github.com/hpcloud/tail"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// TODO -> find pattern to persist the clients on main() scope despite their different typing
type ConnectorInterface interface {
	Open()
	Send(line *tail.Line) bool
	Close()
}

// Connectors factory using ConnectorInterface
func NewConnector(cfg config.Connector) ConnectorInterface {

	switch connType := cfg.Type; connType {
	case "s3":
		return S3Connector{cfg: cfg, client: SetupS3Client(cfg)}
	case "rollbar":
		return RollbarConnector{cfg: cfg}
	case "kafka":
		return KafkaConnector{cfg: cfg}
	default:
		panic("Invalid connector type received in NewConnector function, should have been caught at configuration parsing !")
	}
}
