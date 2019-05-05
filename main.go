package main

import (
	"fmt"
	"os"

	"github.com/moshloop/configadm/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var root = &cobra.Command{
		Use: "configadm",
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
	root.AddCommand(&cmd.CloudInit, &cmd.Minify, &cmd.Apply)
	if len(commit) > 8 {
		version = fmt.Sprintf("%v, commit %v, built at %v", version, commit[0:8], date)
	}
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of configadm",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})

	root.SetUsageTemplate(root.UsageTemplate() + fmt.Sprintf("\nversion: %s\n ", version))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}

}
