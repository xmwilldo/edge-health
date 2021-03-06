package server

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/serf/serf"
	log "k8s.io/klog"
	"k8s.io/klog/v2"

	"github.com/xmwilldo/edge-health/pkg/util"
)

func Server(ctx context.Context, wg *sync.WaitGroup, serf *serf.Serf) {
	srv := &http.Server{Addr: ":" + "8888"}
	http.HandleFunc("/localinfo", func(w http.ResponseWriter, r *http.Request) {
		members := serf.Members()
		data, err := json.Marshal(members)
		if err != nil {
			klog.Errorf("Marshal err: %v", err)
			return
		}
		w.Write(data)
	})

	http.HandleFunc("/debug/flags/v", util.UpdateLogLevel)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server: exit with error: %v", err)
		}
	}()

	for range ctx.Done() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Errorf("Server: program exit, server exit")
		}
		wg.Done()
	}
}
