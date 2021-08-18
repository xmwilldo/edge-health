package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/xmwilldo/edge-health/cmd/app-health/app"
	"github.com/xmwilldo/edge-health/pkg/app-health-daemon/util"

	"k8s.io/component-base/logs"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	ctx, _ := util.SignalWatch()

	command := app.NewAppHealthCommand(ctx)

	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
