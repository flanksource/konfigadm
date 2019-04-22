package cmd

import (
	// . "github.com/moshloop/cloud-config/pkg/cloudmeta"
	. "github.com/moshloop/cloud-config/pkg/cloud-meta"
	"github.com/moshloop/cloud-config/pkg/systemd"

	// "github.com/moshloop/cloud-config/pkg/systemd"
	"github.com/spf13/cobra"
)

var (
	CloudInit = cobra.Command{
		Use:   "cloud-init",
		Short: "Print the version of cloud-config",
		Args:  cobra.MinimumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {
			cfg := SystemConfig{}
			cfg.Init()
			svc := systemd.DefaultSystemdService("kubelet")
			svc.Service.ExecStart = "/bin/echo"

			cfg.Services["kubelet"] = Service{
				Extra: svc,
			}

			cfg.Extra.FinalMessage = "Hwllow world"

			println(cfg.ToCloudInit().String())

		},
	}
)

func init() {
	CloudInit.Flags().String("iso", "", "Create an ISO with the cloud-init embedded as user-metadata")
}
