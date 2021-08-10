package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/xmwilldo/edge-service-autonomy/cmd/webhook/app/options"
	"github.com/xmwilldo/edge-service-autonomy/pkg/util"
	"github.com/xmwilldo/edge-service-autonomy/pkg/webhook"
)

func NewWebhookCommand() *cobra.Command {
	o := options.NewWebHookOptions()
	cmd := &cobra.Command{
		Use: "webhook",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintFlags(cmd.Flags())

			webhookServer, err := webhook.NewWebhookServer(o)
			if err != nil {
				panic(err)
			}
			go webhookServer.Start()

			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
			<-signalChan

			webhookServer.Stop()
		},
	}

	o.AddFlags(cmd.Flags())
	return cmd
}
