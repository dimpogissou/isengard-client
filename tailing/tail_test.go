package tailing

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/hpcloud/tail"
)

var testFile = "test_file.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Creates an empty file to write into
func createEmptyFile(name string) {
	d := []byte("")
	check(ioutil.WriteFile(name, d, 0644))
}

// Test setup function
func testSetup(dir string) *os.File {

	// Create directory
	err := os.Mkdir(dir, 0755)
	check(err)

	// Create test files
	emptyFile, err := os.Create(fmt.Sprintf("%s/%s", dir, testFile))
	check(err)

	return emptyFile
}

// Test teardown function
func testTeardown(dir string) {
	// Delete directory and files
	os.RemoveAll(dir)
}

// Core test method of tailing functionality. Creates a test directory, starts tailing it and asserts generated log lines are correctly received
func TestTailDirectory(t *testing.T) {

	done := make(chan bool)

	timeoutSeconds := 6
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)
	var testDir = "./config_test_files"
	var testLines = 5

	var file = testSetup(testDir)
	defer testTeardown(testDir)

	var testLogLine = "[2020-10-07 20:56:47.375586 UTC][INFO][009] Log message"

	sig := make(chan os.Signal)
	defer close(sig)

	logs := make(chan *tail.Line)

	tails := InitTailsFromDir(testDir)
	for _, t := range tails {
		go SendLines(t, logs)
	}

	func() {
		time.Sleep(1 * time.Second)
		for _ = range make([]int, testLines) {
			file.WriteString(testLogLine)
			file.WriteString("\n")
		}
		file.Close()
	}()

	go func() {
		i := 0
		for line := range logs {
			i += 1
			if line.Text != testLogLine {
				t.Errorf("Log line tailing failed, got [%v], want [%v]", line.Text, testLogLine)
			}
			if i == testLines {
				done <- true
			}
		}
	}()

	select {
	case <-timeout:
		t.Fatalf("Test timed out after %v seconds", timeoutSeconds)
	case <-done:
	}

}
