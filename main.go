package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimpogissou/isengard-server/connectors"
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/dimpogissou/isengard-server/tailing"
	"github.com/hpcloud/tail"
)

func main() {

	// Initialize logger
	logger.InitLogger()

	// Parse CLI argument pointing to configuration file
	configPtr := flag.String("config", "", "Text to parse. (Required)")
	flag.Parse()
	if *configPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Validate and loads config
	cfg := connectors.ValidateAndLoadConfig(configPtr)

	// Create logs channel receiving all tailed loglines from the target directory
	logsChannel := make(chan *tail.Line)

	// Create signal channel listening to interrupt and termination signals
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	// Create and start all connectors, and defer teardown operations
	// TODO -> Also listen for sigChannel here and return so deferred functions are executed on interrupt
	conns := connectors.GenerateConnectors(cfg)
	for _, c := range conns {
		defer c.Close()
	}

	// Tails all log files in directory and sends data to configured targets
	for line := range tailing.TailDirectory(cfg.Directory, logsChannel, sigChannel) {
		for _, conn := range conns {
			go func(c connectors.ConnectorInterface) {
				err := c.Send(line)
				if err != nil {
					// TODO -> Implement fallback logic
					logger.Error("SendError", err.Error())
				}
			}(conn)
		}
	}
}
