package tailing

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var testFile = "test_file.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createEmptyFile(name string) {
	d := []byte("")
	check(ioutil.WriteFile(name, d, 0644))
}

func testSetup(dir string) *os.File {

	// Create directory
	err := os.Mkdir(dir, 0755)
	check(err)

	// Create test files
	emptyFile, err := os.Create(fmt.Sprintf("%s/%s", dir, testFile))
	check(err)

	return emptyFile
}

func testTeardown(dir string) {
	// Delete directory and files
	os.RemoveAll(dir)
}

func TestEndToEnd(t *testing.T) {

	var testDir = "./config_test_files"
	defer testTeardown(testDir)

}
