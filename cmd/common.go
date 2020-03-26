package cmd

import (
	"os"

	_ "github.com/flanksource/konfigadm/pkg" // nolint: golint, stylecheck
	"github.com/flanksource/konfigadm/pkg/phases"
	"github.com/flanksource/konfigadm/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetConfigWithImage(cmd *cobra.Command, args []string, image Image) *types.Config {
	return nil
}

func GetConfig(cmd *cobra.Command, args []string) *types.Config {
	configs, err := cmd.Flags().GetStringSlice("config")
	if err != nil {
		log.Fatalf("%s", err)
	}
	vars, err := cmd.Flags().GetStringSlice("var")
	if err != nil {
		log.Fatalf("%s", err)
	}

	flags := []types.Flag{}

	if alias != nil {
		flags = append(flags, alias.Tags...)
	} else if ok, _ := cmd.Flags().GetBool("detect-tags"); ok {
		for _, _os := range phases.SupportedOperatingSystems {
			if _os.DetectAtRuntime() {
				flags = append(flags, _os.GetTags()...)
			}
		}
		if os.Getenv("container") != "" {
			flags = append(flags, types.CONTAINER)
		}
		log.Debugf("Detected %s\n", flags)
	}

	flagNames, err := cmd.Flags().GetStringSlice("tag")
	if err != nil {
		log.Fatalf("%s", err)
	}

	for _, name := range flagNames {
		if flag, ok := types.FLAG_MAP[name]; ok {
			flags = append(flags, flag)
		} else {
			log.Fatalf("Unknown flag %s", name)
		}
	}

	log.Infof("Using tags: %s\n", flags)

	cfg, err := types.NewConfig(append(configs, args...)...).
		WithVars(vars...).
		WithFlags(flags...).
		Build()

	if err != nil {
		panic(nil)
	}
	return cfg
}
