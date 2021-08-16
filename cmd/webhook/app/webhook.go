package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	ctrwebhook "sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/xmwilldo/edge-service-autonomy/cmd/webhook/app/options"
	"github.com/xmwilldo/edge-service-autonomy/pkg/util"
	"github.com/xmwilldo/edge-service-autonomy/pkg/version"
	"github.com/xmwilldo/edge-service-autonomy/pkg/webhook"
)

func NewWebhookCommand() *cobra.Command {
	o := options.NewDefaultWebHookOptions()
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Start a edge service autonomy webhook server",
		Long:  "Start a edge service autonomy webhook server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(os.Stdout, "edge service autonomy webhook version: %s\n", fmt.Sprintf("%#v", version.Get()))
			if o.VerFlag {
				os.Exit(0)
			}
			util.PrintFlags(cmd.Flags())

			if err := Run(o); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		},
	}
	o.AddFlags(cmd.Flags())
	return cmd
}

func Run(options options.WebHookOptions) error {
	config, err := clientcmd.BuildConfigFromFlags(options.MasterURL, options.KubeConfig)
	if err != nil {
		klog.Fatalf("error setting up webhook's config: %s", err)
	}
	mgr, err := manager.New(config, manager.Options{
		Port:    options.Port,
		CertDir: options.CertDir,
	})
	if err != nil {
		klog.Fatalf("error setting up webhook manager: %s", err)
	}
	hookServer := mgr.GetWebhookServer()

	edgeServiceAutonomy, err := webhook.NewEdgeServiceAutonomy(options)
	if err != nil {
		klog.Fatalf("error setting up webhook manager: %s", err)
	}
	hookServer.Register("/mutating", &ctrwebhook.Admission{Handler: edgeServiceAutonomy})

	hookServer.WebhookMux.Handle("/readyz/", http.StripPrefix("/readyz/", &healthz.Handler{}))

	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		klog.Fatalf("unable to run manager: %s", err)
	}

	return nil
}
