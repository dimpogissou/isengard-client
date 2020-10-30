package testutils

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/dimpogissou/isengard-server/observer"
	"github.com/hpcloud/tail"
)

// Mock Connector implementing ConnectorInterface
type MockConnector struct{}

func (c MockConnector) GetName() string         { return "mockConnector" }
func (c MockConnector) Send(t *tail.Line) error { return nil }
func (c MockConnector) Close() error            { return nil }

// Create test file
func CreateTestFile(dir string, fileName string) *os.File {

	emptyFile, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))
	logger.CheckErrAndPanic(err, "FailedCreatingTestFile", "Unable to create test file")

	return emptyFile
}

// Sleeps then writes to provided file
func SleepThenWriteToFile(file *os.File, duration time.Duration, nLines int, testLine string) {
	time.Sleep(duration)
	for _ = range make([]int, nLines) {
		file.WriteString(testLine)
		file.WriteString("\n")
	}
	file.Close()
}

// Reads lines from subscriber channel, asserts correct number, then sends true to bool channel.
// Should be used with timeout as it will just hang if not enough records are received.
func ReadAndAssertLines(t *testing.T, subscriber observer.Subscriber, logLine string, nLines int, done chan bool) {
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
