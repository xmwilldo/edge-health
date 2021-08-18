package main

import (
	"fmt"
	"os"

	"github.com/xmwilldo/edge-health/cmd/webhook/app"

	"k8s.io/component-base/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	command := app.NewWebhookCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
