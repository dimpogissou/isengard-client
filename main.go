package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimpogissou/isengard-server/connectors"
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/dimpogissou/isengard-server/signals"
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
	signal.Notify(sigChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Create and start all connectors, and defer teardown operations
	// TODO -> Also listen for sigChannel here and return so deferred functions are executed on interrupt
	conns := connectors.CreateConnectors(cfg)

	// Create initial tails
	tails := tailing.InitTailsFromDir(cfg.Directory)

	// Listen to sigChannel and close all connectors and tails if received
	go signals.CloseResourcesOnTerm(sigChannel, logsChannel, conns, tails)

	// Tail existing files
	for _, t := range tails {
		go tailing.SendLines(t, logsChannel)
	}

	// Routine reading logs lines and sending them to configured connectors
	connectors.SendToConnectors(logsChannel, conns)
}
