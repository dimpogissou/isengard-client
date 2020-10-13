package tailing

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hpcloud/tail"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("main")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func TailDirectory(dir string, logsChannel chan *tail.Line, sigChannel chan os.Signal) chan *tail.Line {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error(fmt.Sprintf("Could not get files from directory %s due to -> %s", dir, err))
	}

	var tails = make([]*tail.Tail, 0)
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dir, file.Name())

		log.Info(fmt.Sprintf("Start tailing file %s", filePath))

		t, err := tail.TailFile(filePath, tail.Config{Follow: true, MustExist: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}, ReOpen: true, Poll: true})

		if err != nil {
			log.Error(fmt.Sprintf("Could not tail file [%s] due to -> %s", filePath, err))
		}

		tails = append(tails, t)
	}

	go func() {
		<-sigChannel
		for i, t := range tails {
			log.Info("Closing tail channel for file ", i)
			t.Stop()
		}
		log.Info("Closing logsChannel for dir ", dir)
		close(logsChannel)
	}()

	for _, t := range tails {
		go func(t *tail.Tail) {
			for line := range t.Lines {
				logsChannel <- line
			}
			log.Debug("Closing tailing goroutine")
		}(t)
	}

	return logsChannel
}
