package connectors

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/hpcloud/tail"
	"gopkg.in/fsnotify.v1"
)

const timeoutSeconds = 3

const testDir = "./config_test_files"

const testLogLine = "[2020-10-07 20:56:47.375586 UTC][INFO][009] Log message"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Test setup function
func createTestFile(dir string, fileName string) *os.File {

	emptyFile, err := os.Create(fmt.Sprintf("%s/%s", dir, fileName))
	check(err)

	return emptyFile
}

// Test teardown function
func testTeardown(dir string) {
	// Delete directory and files
	os.RemoveAll(dir)
}

type mockConnector struct{}

func (c mockConnector) Send(t *tail.Line) error { return nil }
func (c mockConnector) Close() error            { return nil }

func sleepThenWriteToFile(file *os.File, duration time.Duration, nLines int, testLine string) {
	time.Sleep(duration)
	for _ = range make([]int, nLines) {
		file.WriteString(testLine)
		file.WriteString("\n")
	}
	file.Close()
}

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

// Core test method of tailing functionality. Creates a test directory, starts tailing it and asserts generated log lines are correctly received
func TestTailDirectory(t *testing.T) {

	done := make(chan bool)
	defer close(done)
	const fileName = "test_file_1.txt"
	const nLines = 5
	var timeout = time.After(time.Duration(timeoutSeconds) * time.Second)

	// Create test directory
	err := os.Mkdir(testDir, 0755)
	check(err)
	defer testTeardown(testDir)

	// Create test file
	testFile1 := createTestFile(testDir, fileName)

	// Create signal channel
	sigCh := make(chan os.Signal)
	defer close(sigCh)

	// Create publisher/subscriber with mockConnector
	logsPublisher := Publisher{}
	logsCh := make(chan *tail.Line)
	subscriber := Subscriber{
		Channel:   logsCh,
		Connector: mockConnector{},
	}
	logsPublisher.Subscribe(subscriber.Channel)

	// Create tail goroutines
	tails := InitTailsFromDir(testDir)
	for _, t := range tails {
		go TailAndPublish(t.Lines, logsPublisher)
	}

	// Assert tailing of existing files works correctly
	go sleepThenWriteToFile(testFile1, 1*time.Second, nLines, testLogLine)
	go readAndAssertLines(t, subscriber, testLogLine, nLines, done)

	select {
	case <-timeout:
		t.Fatalf("Test timed out after %v seconds", timeoutSeconds)
	case <-done:
	}

}

func TestTailNewFiles(t *testing.T) {

	done := make(chan bool)
	defer close(done)
	const fileName = "test_file_2.txt"
	const nLines = 5
	var timeout = time.After(time.Duration(timeoutSeconds) * time.Second)

	// Create test directory
	err := os.Mkdir(testDir, 0755)
	check(err)
	defer testTeardown(testDir)

	// Create FS events watcher detecting new files
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	err = watcher.Add(testDir)
	if err != nil {
		log.Fatal(err)
	}

	// Create publisher/subscriber with mockConnector
	logsPublisher := Publisher{}
	logsCh := make(chan *tail.Line)
	subscriber := Subscriber{
		Channel:   logsCh,
		Connector: mockConnector{},
	}
	logsPublisher.Subscribe(subscriber.Channel)

	// Create signal channel
	sigCh := make(chan os.Signal)
	defer close(sigCh)

	// Add new file and ensure watcher picks it up and starts tailing it from start
	go TailNewFiles(watcher, logsPublisher, sigCh)
	testFile2 := createTestFile(testDir, fileName)
	go sleepThenWriteToFile(testFile2, 1*time.Second, nLines, testLogLine)
	go readAndAssertLines(t, subscriber, testLogLine, nLines, done)

	select {
	case <-timeout:
		t.Fatalf("Test timed out after %v seconds", timeoutSeconds)
	case <-done:
	}

}
