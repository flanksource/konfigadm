package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Minify = cobra.Command{
		Use:   "minify",
		Short: "Resolve all lookups and dependencies and export a single config file",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)
