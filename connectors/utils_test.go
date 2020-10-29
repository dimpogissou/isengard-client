package connectors

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hpcloud/tail"
)

// Mock Connector implementing ConnectorInterface
type mockConnector struct{}

func (c mockConnector) GetName() string         { return "mockConnector" }
func (c mockConnector) Send(t *tail.Line) error { return nil }
func (c mockConnector) Close() error            { return nil }

// Create test file
func createTestFile(dir string, fileName string) *os.File {

	emptyFile, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))
	check(err)

	return emptyFile
}

// Sleeps then writes to provided file
func sleepThenWriteToFile(file *os.File, duration time.Duration, nLines int, testLine string) {
	time.Sleep(duration)
	for _ = range make([]int, nLines) {
		file.WriteString(testLine)
		file.WriteString("\n")
	}
	file.Close()
}

// Reads lines from subscriber channel, asserts correct number, then sends true to bool channel.
// Should be used with timeout as it will just hang if not enough records are received.
func readAndAssertLines(t *testing.T, subscriber Subscriber, logLine string, nLines int, done chan bool) {
	i := 0
	for line := range subscriber.Channel {
		i += 1
		if line.Text != logLine {
			t.Errorf("Log line tailing failed, got [%v], want [%v]", line.Text, logLine)
		}
		if i == nLines {
			done <- true
		}
	}
}
