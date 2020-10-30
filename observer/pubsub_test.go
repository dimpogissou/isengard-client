package observer_test

import (
	"testing"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/dimpogissou/isengard-server/observer"
	"github.com/dimpogissou/isengard-server/testutils"
	"github.com/hpcloud/tail"
)

func TestPubSub(t *testing.T) {

	const timeoutSeconds = 3

	const testDir = "./config_test_files"

	const testLogLine = "[2020-10-07 20:56:47.375586 UTC][INFO][009] Log message"

	// Create publisher/subscribers with mockConnector
	logger.Info("Creating publisher and subscribers ...")
	logsPublisher := observer.Publisher{}

	ch1 := make(chan *tail.Line)
	conn1 := testutils.MockConnector{}
	defer close(ch1)
	subscriber1 := observer.Subscriber{
		Channel:   ch1,
		Connector: conn1,
	}

	ch2 := make(chan *tail.Line)
	conn2 := testutils.MockConnector{}
	defer close(ch2)
	subscriber2 := observer.Subscriber{
		Channel:   ch2,
		Connector: conn2,
	}
	logsPublisher.Subscribe(subscriber1.Channel)
	logsPublisher.Subscribe(subscriber2.Channel)

	// Create test log line
	testLine := tail.Line{Text: "logMessage", Err: nil}

	// Write to publisher and assert both subscribers received it
	logger.Info("Publishing log line ...")
	go logsPublisher.Publish(&testLine)

	l1 := <-ch1
	l2 := <-ch2

	if l1.Text != testLine.Text || l2.Text != testLine.Text {
		t.Errorf("Failed asserting log lines send to subscriber channels, got %s and %s, want %s", l1.Text, l2.Text, testLine.Text)
	}

}
