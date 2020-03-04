package cmd

import (
	"io/ioutil"
	"net/http"
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
				content := []byte(file.Content)
				if file.Content == "" && file.ContentFromURL != "" {
					log.Infof("Downloading %s to path %s", file.ContentFromURL, path)
					resp, err := http.Get(file.ContentFromURL)
					if err != nil {
						log.Errorf("Failed to download from url %s: %v", file.ContentFromURL, err)
						continue
					}
					defer resp.Body.Close()
					c, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Errorf("Failed to read response body from url %s: %v", file.ContentFromURL, err)
						continue
					}
					content = c
				}
				ioutil.WriteFile(path, []byte(content), os.FileMode(perms))
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
