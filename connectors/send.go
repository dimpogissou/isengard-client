package connectors

import (
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

// Reads log lines from channel and sends to all configured connectors
func SendToConnectors(ch chan *tail.Line, conns []ConnectorInterface) {

	for line := range ch {
		for _, conn := range conns {
			go func(c ConnectorInterface) {
				err := c.Send(line)
				if err != nil {
					// TODO -> Implement fallback logic
					logger.Error("SendError", err.Error())
				}
			}(conn)
		}
	}
	logger.Info("Logs channel closed, stop sending to connectors")
}
