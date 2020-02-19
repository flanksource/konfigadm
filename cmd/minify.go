package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

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
			primitive, _ := cmd.Flags().GetBool("primitive")
			if primitive {
				data, _ := yaml.Marshal(map[string]interface{}{
					"filesystem": fs,
					"commands":   commands,
				})
				fmt.Println(string(data))
			} else {
				data, _ := yaml.Marshal(cfg)
				fmt.Println(string(data))

			}
		},
	}
)

func init() {
	Minify.Flags().Bool("primitive", false, "Minify down to primitive level of commands and files only")
}
