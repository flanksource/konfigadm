package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"gopkg.in/flanksource/yaml.v3"
)

var primitive bool
var bash bool
var (
	Minify = cobra.Command{
		Use:   "minify",
		Short: "Resolve all lookups and dependencies and export a single config file",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := GetConfig(cmd, args)
			fs, commands, err := cfg.ApplyPhases()
			if err != nil {
				log.Fatalf("Error applying phases %s\n", err)
			}

			if primitive {
				data, _ := yaml.Marshal(map[string]interface{}{
					"filesystem": fs,
					"commands":   commands,
				})
				fmt.Println(string(data))
			} else if bash {
				if out, err := cfg.ToBash(); err != nil {
					log.Fatalf("Error converting to bas: %v", err)
				} else {
					fmt.Println(out)
				}
			} else {
				data, _ := yaml.Marshal(cfg)
				fmt.Println(string(data))

			}
		},
	}
)

func init() {
	Minify.Flags().BoolVar(&primitive, "primitive", false, "Minify down to primitive level of commands and files only")
	Minify.Flags().BoolVar(&bash, "bash", false, "Export a single bash file with base64 encoded files")
}
