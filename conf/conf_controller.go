package conf

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var AgentConfData map[string]string

func init() {
	AgentConfData = make(map[string]string)
	appDir, err := os.Getwd()
	if err != nil {
		panic("Get app dir error 01.")
	}
	confFile := appDir + "/conf/agent.conf"

	f, err := os.Open(confFile)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Close config file error: ", f.Name())
		}
	}(f)
	if err != nil {
		fmt.Println(err.Error()) // change to logger
		return
	}
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") {
			continue
		}
		items := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(items[0])
		value := strings.TrimSpace(items[1])
		AgentConfData[key] = value
	}
}
