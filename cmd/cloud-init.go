package cmd

import (
	. "github.com/moshloop/configadm/pkg/phases"
	"github.com/spf13/cobra"
)

var (
	CloudInit = cobra.Command{
		Use:   "cloud-init",
		Short: "Print the version of cloud-config",
		Args:  cobra.MinimumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			configs, err := cmd.Flags().GetStringSlice("config")
			if err != nil {
				panic(err)
			}
			vars, err := cmd.Flags().GetStringSlice("var")
			cfg, err := NewSystemConfig(vars, configs)

			if err != nil {
				panic(nil)
			}
			println(cfg.ToCloudInit().String())

		},
	}
)

func init() {
	CloudInit.Flags().String("iso", "", "Create an ISO with the cloud-init embedded as user-metadata")
}
