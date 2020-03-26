package main

import (
	"fmt"
	"os"

	"github.com/flanksource/konfigadm/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	var root = &cobra.Command{
		Use: "konfigadm",
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
			log.SetOutput(os.Stderr)
		},
	}

	root.PersistentFlags().StringSliceP("config", "c", []string{}, "Config file path in YAML or JSON format, Use - to read YAML from stdin")
	root.PersistentFlags().StringSliceP("var", "e", []string{}, "Extra Variables to in key=value format ")
	root.PersistentFlags().StringSliceP("tag", "t", []string{}, "Runtime tags to use, valid tags:  debian,ubuntu,redhat,rhel,centos,aws,vmware")
	root.PersistentFlags().BoolP("detect-tags", "d", true, "Detect tags to use")

	root.AddCommand(&cmd.CloudInit, &cmd.Minify, &cmd.Apply, &cmd.Verify, &cmd.Images)
	if len(commit) > 8 {
		version = fmt.Sprintf("%v, commit %v, built at %v", version, commit[0:8], date)
	}
	root.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version of konfigadm",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})

	root.PersistentFlags().CountP("loglevel", "v", "Increase logging level")
	root.SetUsageTemplate(root.UsageTemplate() + fmt.Sprintf("\nversion: %s\n ", version))

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
