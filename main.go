package main

import (
	"fmt"
	"opstack_agent/conf"
	"opstack_agent/logs"
)

func main() {
	fmt.Println(conf.AgentConfData)
	logs.Logger.Debug("test line1")
}
