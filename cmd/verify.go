package cmd

import (
	"os"

	"github.com/moshloop/configadm/pkg/types"
	"github.com/spf13/cobra"
)

var (

	//Verify command
	Verify = cobra.Command{
		Use:   "verify",
		Short: "Verify that the configuration has been applied correctly and is in a healthy state",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := GetConfig(cmd)
			cfg.ApplyPhases()
			verifier := types.VerifyResults{}
			if !cfg.Verify(&verifier) {
				verifier.Done()
				os.Exit(1)
			}
			verifier.Done()
		},
	}
)
