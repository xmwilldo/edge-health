package app

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/xmwilldo/edge-health/cmd/app-health/app/options"
	"github.com/xmwilldo/edge-health/pkg/app-health-daemon/daemon"
	"github.com/xmwilldo/edge-health/pkg/util"
)

func NewAppHealthCommand(ctx context.Context) *cobra.Command {
	o := options.NewAppHealthOptions()
	cmd := &cobra.Command{
		Use: "app-health",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrintFlags(cmd.Flags())

			d := daemon.NewAppDaemon(o)
			d.Run(ctx)
		},
	}

	o.AddFlags(cmd.Flags())

	return cmd
}
