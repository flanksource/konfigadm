package cmd

import (
	_ "github.com/moshloop/konfigadm/pkg"
	"github.com/moshloop/konfigadm/pkg/phases"
	"github.com/moshloop/konfigadm/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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

	if ok, _ := cmd.Flags().GetBool("detect-tags"); ok {
		for _, _os := range phases.SupportedOperatingSystems {
			if _os.DetectAtRuntime() {
				log.Infof("Detected %s\n", _os.GetTags())
				flags = append(flags, _os.GetTags()...)
			}
		}
	}

	flagNames, err := cmd.Flags().GetStringSlice("tag")
	for _, name := range flagNames {

		if flag, ok := types.FLAG_MAP[name]; ok {
			flags = append(flags, flag)
		} else {
			log.Fatalf("Unknown flag %s", name)
		}

	}

	log.Debugf("Using tags: %s\n", flags)
	if err != nil {
		log.Fatalf("%s", err)
	}

	cfg, err := types.NewConfig(append(configs, args...)...).
		WithVars(vars...).
		WithFlags(flags...).
		Build()

	if err != nil {
		panic(nil)
	}
	return cfg

}
