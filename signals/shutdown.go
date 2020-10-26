package signals

import (
	"fmt"
	"os"

	"github.com/dimpogissou/isengard-server/connectors"
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

// Closes logs channel, connectors, and tails on interruption signals
func CloseResourcesOnTerm(sigCh chan os.Signal, logsCh chan *tail.Line, connectors []connectors.ConnectorInterface, tails []*tail.Tail) {
	<-sigCh
	for _, c := range connectors {
		// TODO -> add GetName to connector interface and pass it here
		err := c.Close()
		if err != nil {
			logger.Error("FailedClosingConnector", fmt.Sprintf("Failed closing connector: %s", err.Error()))
		}
	}
	for _, t := range tails {
		logger.Info("Closing tail channel")
		t.Stop()
	}
	logger.Info("Closing logs channel")
	close(logsCh)
}
