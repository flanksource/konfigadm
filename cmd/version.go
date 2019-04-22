package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	VersionInfo = "develop"

	Version = cobra.Command{
		Use:   "version",
		Short: "Print the version of cloud-config",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(VersionInfo)
		},
	}
)
