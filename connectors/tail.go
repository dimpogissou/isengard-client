package connectors

import (
	"fmt"
	"io/ioutil"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

// Collects file names in provided directory as an array of strings
func getFileNamesInDir(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Error("FailedRetrievingFiles", fmt.Sprintf("Could not get files from directory %s due to -> %s", dir, err))
	}
	paths := []string{}
	for _, f := range files {
		paths = append(paths, f.Name())
	}
	return paths
}

// Starts tailing a file at provided path
func createTail(path string) (*tail.Tail, error) {

	logger.Info(fmt.Sprintf("Start tailing file %s", path))

	t, err := tail.TailFile(path, tail.Config{Follow: true, MustExist: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}, ReOpen: true, Poll: true})

	return t, err
}

// Starts tailing all files in provided directory
func InitTailsFromDir(dir string) []*tail.Tail {

	var files = getFileNamesInDir(dir)
	var tails = make([]*tail.Tail, 0)
	for _, fileName := range files {
		filePath := fmt.Sprintf("%s/%s", dir, fileName)
		t, err := createTail(filePath)
		if err != nil {
			logger.Error("FailedTailingFile", fmt.Sprintf("Could not tail file [%s] due to -> %s", filePath, err))
		} else {
			tails = append(tails, t)
		}
	}
	return tails
}

// Routine tailing a file and sending lines to logs channel
func TailAndPublish(lines chan *tail.Line, publisher Publisher) {
	for line := range lines {
		publisher.Publish(line)
	}
}
