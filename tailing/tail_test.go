package tailing

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/dimpogissou/isengard-server/observer"
	"github.com/dimpogissou/isengard-server/testutils"
	"github.com/hpcloud/tail"
	"gopkg.in/fsnotify.v1"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Cleanup test dir
func testTeardown(dir string) {
	// Delete directory and files
	os.RemoveAll(dir)
}

// Creates a test directory and file, starts tailing it and asserts generated log lines are correctly received
func TestTailDirectory(t *testing.T) {

	const timeoutSeconds = 3

	const testDir = "./config_test_files"

	const testLogLine = "[2020-10-07 20:56:47.375586 UTC][INFO][009] Log message"

	// Test specific const/vars
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
	testFile1 := testutils.CreateTestFile(testDir, fileName)

	// Create signal channel
	sigCh := make(chan os.Signal)
	defer close(sigCh)

	// Create publisher/subscriber with mockConnector
	logsPublisher := observer.Publisher{}
	logsCh := make(chan *tail.Line)
	subscriber := observer.Subscriber{
		Channel:   logsCh,
		Connector: testutils.MockConnector{},
	}
	logsPublisher.Subscribe(subscriber.Channel)

	// Create tail goroutines
	tails := InitTailsFromDir(testDir)
	for _, t := range tails {
		go TailAndPublish(t.Lines, logsPublisher)
		defer t.Stop()
	}

	// Assert tailing of existing files works correctly
	go testutils.SleepThenWriteToFile(testFile1, 1*time.Second, nLines, testLogLine)
	go testutils.ReadAndAssertLines(t, subscriber, testLogLine, nLines, done)

	select {
	case <-timeout:
		t.Fatalf("Test timed out after %v seconds", timeoutSeconds)
	case <-done:
	}

}

// Creates a test directory, starts watcher, creates and writes to a new file and asserts generated log lines are correctly received
func TestTailNewFiles(t *testing.T) {

	const timeoutSeconds = 3

	const testDir = "./config_test_files"

	const testLogLine = "[2020-10-07 20:56:47.375586 UTC][INFO][009] Log message"

	// Test specific const/vars
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
	check(err)
	defer watcher.Close()
	err = watcher.Add(testDir)
	check(err)

	// Create publisher/subscriber with mockConnector
	logsPublisher := observer.Publisher{}
	logsCh := make(chan *tail.Line)
	subscriber := observer.Subscriber{
		Channel:   logsCh,
		Connector: testutils.MockConnector{},
	}
	logsPublisher.Subscribe(subscriber.Channel)

	// Create signal channel
	sigCh := make(chan os.Signal)
	defer close(sigCh)

	// Add new file and ensure watcher picks it up and starts tailing it from start
	go TailNewFiles(watcher, logsPublisher, sigCh)
	testFile2 := testutils.CreateTestFile(testDir, fileName)
	testutils.SleepThenWriteToFile(testFile2, 1*time.Second, nLines, testLogLine)
	go testutils.ReadAndAssertLines(t, subscriber, testLogLine, nLines, done)

	select {
	case <-timeout:
		t.Fatalf("Test timed out after %v seconds", timeoutSeconds)
	case <-done:
		sigCh <- syscall.SIGINT
	}

}
