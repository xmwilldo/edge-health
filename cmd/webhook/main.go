package main

import (
	"fmt"
	"os"

	"k8s.io/component-base/logs"

	"github.com/xmwilldo/edge-service-autonomy/cmd/webhook/app"
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
