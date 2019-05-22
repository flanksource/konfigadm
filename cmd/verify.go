package cmd

import (
	"os"

	"github.com/moshloop/konfigadm/pkg/types"
	log "github.com/sirupsen/logrus"
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
			_, _, err := cfg.ApplyPhases()
			if err != nil {
				log.Error(err)
			}
			verifier := types.VerifyResults{}
			if !cfg.Verify(&verifier) {
				verifier.Done()
				os.Exit(1)
			}
			verifier.Done()
		},
	}
)
