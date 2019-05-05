package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var (
	Minify = cobra.Command{
		Use:   "minify",
		Short: "Resolve all lookups and dependencies and export a single config file",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := GetConfig(cmd)
			cfg.ApplyPhases()
			data, _ := yaml.Marshal(cfg)
			fmt.Println(string(data))
		},
	}
)
