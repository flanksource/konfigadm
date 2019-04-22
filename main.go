package main

import (
	"os"

	"github.com/moshloop/cloud-config/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {

	var root = &cobra.Command{
		Use: "cloud-config",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level, _ := cmd.Flags().GetCount("loglevel")
			switch {
			case level > 1:
				log.SetLevel(log.DebugLevel)
			case level > 0:
				log.SetLevel(log.InfoLevel)
			default:
				log.SetLevel(log.WarnLevel)
			}
		},
	}

	// root.PersistentFlags().StringP("inventory", "i", "", "Specify inventory host path or comma separated host list")
	// root.PersistentFlags().Bool("version", false, "")
	// root.PersistentFlags().StringSliceP("extra-vars", "e", []string{}, "Set additional variables as key=value or YAML/JSON, if filename prepend with @")
	// root.PersistentFlags().StringP("limit", "l", "", "Limit selected hosts to an additional pattern")
	root.PersistentFlags().CountP("loglevel", "v", "Increase logging level")
	root.AddCommand(&cmd.Version, &cmd.CloudInit)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}

}
