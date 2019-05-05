package cmd

import (
	_ "github.com/moshloop/configadm/pkg"
	"github.com/spf13/cobra"
)

var (
	CloudInit = cobra.Command{
		Use:   "cloud-init",
		Short: "Print the version of cloud-config",
		Args:  cobra.MinimumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			cfg := GetConfig(cmd)
			println(cfg.ToCloudInit().String())

		},
	}
)

func init() {
	CloudInit.Flags().String("iso", "", "Create an ISO with the cloud-init embedded as user-metadata")
}
