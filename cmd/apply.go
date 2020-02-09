package cmd

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Apply = cobra.Command{
		Use:   "apply",
		Short: "Apply the configuration to the local machine",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cfg := GetConfig(cmd, args)

			files, commands, err := cfg.ApplyPhases()
			if err != nil {
				panic(err)
			}

			for path, file := range files {
				log.Infof("Writing %s\n", utils.LightGreenf(path))
				perms, _ := strconv.Atoi(file.Permissions)
				if perms == 0 {
					perms = 0644
				}
				ioutil.WriteFile(path, []byte(file.Content), os.FileMode(perms))
			}

			for _, cmd := range commands {
				log.Infof("Executing %s\n", utils.LightGreenf(cmd.Cmd))
				if err := utils.Exec(cmd.Cmd); err != nil {
					log.Fatalf("Failed to run: %s, %s", cmd.Cmd, err)
				}
			}

		},
	}
)
