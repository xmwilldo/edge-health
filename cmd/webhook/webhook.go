package main

import (
	"os"

	"github.com/xmwilldo/edge-service-autonomy/cmd/webhook/app"
)

func main() {
	command := app.NewWebhookCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
