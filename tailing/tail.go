package tailing

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dimpogissou/isengard-server/logger"
	"github.com/hpcloud/tail"
)

// Core function tailing all files in a directory and sending them to the logs channel
// Listens to signal channel for interruptions and closes tailing processes before teardown
// Returns the logsChannel so it can be used as an interator
func TailDirectory(dir string, logsChannel chan *tail.Line, sigChannel chan os.Signal) chan *tail.Line {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		logger.Error("FailedRetrievingFiles", fmt.Sprintf("Could not get files from directory %s due to -> %s", dir, err))
	}

	var tails = make([]*tail.Tail, 0)
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dir, file.Name())

		logger.Info(fmt.Sprintf("Start tailing file %s", filePath))

		t, err := tail.TailFile(filePath, tail.Config{Follow: true, MustExist: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}, ReOpen: true, Poll: true})

		if err != nil {
			logger.Error("FailedTailingFile", fmt.Sprintf("Could not tail file [%s] due to -> %s", filePath, err))
		}

		tails = append(tails, t)
	}

	go func() {
		<-sigChannel
		for i, t := range tails {
			logger.Info(fmt.Sprintf("Closing tail channel for file %v", i))
			t.Stop()
		}
		logger.Info(fmt.Sprintf("Closing logsChannel for directory %s", dir))
		close(logsChannel)
	}()

	for _, t := range tails {
		go func(t *tail.Tail) {
			for line := range t.Lines {
				logsChannel <- line
			}
			logger.Debug("Closing tailing goroutine")
		}(t)
	}

	return logsChannel
}
