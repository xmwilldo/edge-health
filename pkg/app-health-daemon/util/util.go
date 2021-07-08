package util

import (
	"context"
	"os"
	"os/signal"

	"golang.org/x/sys/unix"
)

func SignalWatch() (context.Context, context.CancelFunc) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, unix.SIGTERM, unix.SIGINT)
	go func() {
		for range signals {
			cancelFunc()
		}
	}()
	return ctx, cancelFunc
}
