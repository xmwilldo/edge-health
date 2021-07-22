package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/serf/serf"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog"

	"github.com/xmwilldo/edge-service-autonomy/cmd/app-health-daemon/app/options"
	"github.com/xmwilldo/edge-service-autonomy/pkg/app-health-daemon/server"
)

type AppDaemon struct {
	serf      *serf.Serf
	joinIp    string
	namespace string
	svcName   string
}

func NewAppDaemon(options *options.AppHealthOptions) *AppDaemon {

	config := serf.DefaultConfig()
	serf, err := serf.Create(config)
	if err != nil {
		return nil
	}

	return &AppDaemon{
		serf:      serf,
		joinIp:    options.JoinIp,
		namespace: options.Namespace,
		svcName:   options.SvcName,
	}
}

func (d *AppDaemon) Run(ctx context.Context) {
	wg := sync.WaitGroup{}

	err := d.serf.SetTags(map[string]string{"namespace": d.namespace, "svc_name": d.svcName})
	if err != nil {
		klog.Errorf("set tags err: %v", err)
		return
	}

	if d.joinIp != "" {
		_, err := d.serf.Join([]string{d.joinIp}, true)
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
