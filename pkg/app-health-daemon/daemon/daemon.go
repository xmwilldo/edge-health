package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/serf/serf"
	"github.com/xmwilldo/edge-service-autonomy/cmd/app-health-daemon/app/options"
	"github.com/xmwilldo/edge-service-autonomy/pkg/app-health-daemon/server"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"
)

type AppDaemon struct {
	serf   *serf.Serf
	joinIP string
}

func NewAppDaemon(options *options.AppHealthOptions) *AppDaemon {

	config := serf.DefaultConfig()
	serf, err := serf.Create(config)
	if err != nil {
		return nil
	}

	return &AppDaemon{
		serf:   serf,
		joinIP: options.JoinIP,
	}
}

func (d *AppDaemon) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	if d.joinIP != "" {
		_, err := d.serf.Join([]string{d.joinIP}, false)
		if err != nil {
			klog.Errorf("join serf err: %v", err)
			return
		}
	}

	go wait.Until(d.PrintMembers, time.Second*3, ctx.Done())

	wg.Add(1)
	go server.Server(ctx, &wg, d.serf)

	for range ctx.Done() {
		wg.Wait()
		return
	}
}

func (d *AppDaemon) PrintMembers() {
	members := d.serf.Members()
	fmt.Printf("Members: %+v\n", members)
}
