package main

import (
	"fmt"
	"opstack_agent/conf"
	"opstack_agent/logs"
	"time"
)

func main() {
	fmt.Println(conf.AgentConfData)
	logs.Logger.Debug("test line1")
	time.Sleep(2 * time.Second)
}
