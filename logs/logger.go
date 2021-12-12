package logs

import (
	"fmt"
	"log"
	"opstack_agent/conf"
	"os"
	"time"
)

func getLogPath(dt string) string {
	appDir, err := os.Getwd()
	if err != nil {
		panic("Get app dir error 01.")
	}
	logName := conf.AgentConfData["logName"]
	return fmt.Sprintf("%s/logs/%s_%s.log", appDir, logName, dt)
}

type LogHandler struct {
	logger   *log.Logger
	dtSuffix string
}

func (lh *LogHandler) rotate() {
	year, month, day := time.Now().Date()
	dt := fmt.Sprintf("%d-%d-%d", year, month, day)
	if lh.logger == nil {
		lh.logger = log.New(os.Stdout, "", log.LstdFlags)
	}
	if lh.dtSuffix != dt {
		newLogName := getLogPath(dt)
		newLogFile, err := os.OpenFile(newLogName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic("Open log f error: " + newLogName + err.Error())
		}
		lh.logger.SetOutput(newLogFile)
	}
}

func (lh *LogHandler) write(s, l string) {
	lh.rotate()
	lh.logger.Printf("%s;level=%s\n", s, l)
}

func (lh *LogHandler) Debug(s string) {
	lh.write(s, "debug")
}

var Logger *LogHandler

func init() {
	Logger = new(LogHandler)
	Logger.rotate()
}
