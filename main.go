package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimpogissou/isengard-server/config"
	"github.com/dimpogissou/isengard-server/tailing"
	"github.com/hpcloud/tail"
)

func main() {
	cfg := config.ValidateAndLoadConfig()
	regex := config.BuildRegex(cfg)

	logsChannel := make(chan *tail.Line)
	sigChannel := make(chan os.Signal)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	for line := range tailing.TailDirectory(cfg.Directory, logsChannel, sigChannel) {
		// TODO -> send to connectors
		fmt.Printf("%v", tailing.ParseLine(line, regex))
	}
}
