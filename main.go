package main

import (
	"os"

	"github.com/moshloop/configadm/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {

	var root = &cobra.Command{
		Use: "cloud-config",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level, _ := cmd.Flags().GetCount("loglevel")
			switch {
			case level > 2:
				log.SetLevel(log.TraceLevel)
			case level > 1:
				log.SetLevel(log.DebugLevel)
			case level > 0:
				log.SetLevel(log.InfoLevel)
			default:
				log.SetLevel(log.WarnLevel)
			}
		},
	}

	root.PersistentFlags().StringSliceP("config", "c", []string{}, "Config files in YAML or JSON format")
	root.PersistentFlags().StringSliceP("var", "e", []string{}, "Variables")
	root.PersistentFlags().StringSliceP("tag", "t", []string{}, "Runtime tags to set")
	root.PersistentFlags().CountP("loglevel", "v", "Increase logging level")
	root.AddCommand(&cmd.Version, &cmd.CloudInit, &cmd.Minify)

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}

}
