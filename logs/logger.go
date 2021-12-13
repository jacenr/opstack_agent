package logs

import (
	"fmt"
	"log"
	"opstack_agent/conf"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex

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
		lh.dtSuffix = dt

		// clean up old log files
		dayNum := conf.AgentConfData["logDays"]
		days, _ := strconv.Atoi(dayNum)
		go func(logName string, rotateDay int) {
			lock.Lock()
			defer lock.Unlock()
			err := filepath.Walk(filepath.Dir(logName), func(path string, info os.FileInfo, err error) error {
				if info.IsDir() {
					return nil
				}
				p, _ := regexp.Compile("agent_(\\d{4}-\\d{1,2}-\\d{1,2})\\.log")
				if !p.MatchString(info.Name()) {
					return nil
				}
				dtStr := p.FindStringSubmatch(info.Name())[1]
				dtTmp, err := time.Parse("2006-1-2", dtStr)
				if err != nil {
					return err
				}

				if dtTmp.Before(time.Now().Add(-86400 * time.Duration(rotateDay) * time.Second)) {
					err = os.Remove(path)
					if err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				fmt.Printf("scan %s err: %s\n", filepath.Dir(newLogName), err.Error())
			}
		}(newLogName, days)
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
	lock = sync.Mutex{}
	Logger = new(LogHandler)
	Logger.rotate()
}
