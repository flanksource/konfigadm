package cmd

import (
	_ "github.com/moshloop/configadm/pkg"
	"github.com/moshloop/configadm/pkg/types"
	log "github.com/sirupsen/logrus"
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
				log.Fatalf("%s", err)
			}
			vars, err := cmd.Flags().GetStringSlice("var")
			if err != nil {
				log.Fatalf("%s", err)
			}

			flags := []types.Flag{}
			flagNames, err := cmd.Flags().GetStringSlice("tag")
			for _, name := range flagNames {

				if flag, ok := types.FLAG_MAP[name]; ok {
					flags = append(flags, flag)
				} else {
					log.Fatalf("Unknown flag %s", name)
				}

			}
			if err != nil {
				log.Fatalf("%s", err)
			}

			cfg, err := types.NewConfig(configs...).
				WithVars(vars...).
				WithFlags(flags...).
				Build()

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
