package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimpogissou/isengard-server/connectors"
	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
	"gopkg.in/fsnotify.v1"
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

	// Create signal channel listening to interrupt and termination signals
	sigChannel := make(chan os.Signal)
	defer close(sigChannel)
	signal.Notify(sigChannel, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Create logs publisher
	logsPublisher := connectors.Publisher{}

	// Create FS events watcher detecting new files
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(cfg.Directory)
	if err != nil {
		log.Fatal(err)
	}

	// Start all configured connectors
	conns := connectors.CreateConnectors(cfg)

	// Subscribe to logsPublisher for each connector
	for _, conn := range conns {
		ch := make(chan *tail.Line)
		defer close(ch)
		defer conn.Close()
		subscriber := connectors.Subscriber{
			Channel:   ch,
			Connector: conn,
		}
		logsPublisher.Subscribe(subscriber.Channel)
		go subscriber.ListenToChannel()
	}

	// Publish lines for each file in a separate thread
	tails := connectors.InitTailsFromDir(cfg.Directory)
	for _, t := range tails {
		defer t.Stop()
		go connectors.TailAndPublish(t.Lines, logsPublisher)
	}

	// Watch for new files added and start tailing them, return on interruption signal to execute deferred calls
	connectors.TailNewFiles(watcher, logsPublisher, sigChannel)
}
