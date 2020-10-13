package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimpogissou/isengard-server/config"
	"github.com/dimpogissou/isengard-server/connectors"
	"github.com/dimpogissou/isengard-server/tailing"
	"github.com/hpcloud/tail"
)

func main() {

	configPtr := flag.String("config", "", "Text to parse. (Required)")
	flag.Parse()
	if *configPtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := config.ValidateAndLoadConfig(configPtr)

	logsChannel := make(chan *tail.Line)
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	conns := []connectors.ConnectorInterface{}
	for _, connCfg := range cfg.Connectors {
		c := connectors.NewConnector(connCfg)
		conns = append(conns, c)
		c.Open()
		defer c.Close()
	}

	for line := range tailing.TailDirectory(cfg.Directory, logsChannel, sigChannel) {
		for _, conn := range conns {
			go conn.Send(line)
		}
	}
}
