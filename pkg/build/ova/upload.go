package ova

import (
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/flanksource/konfigadm/pkg/utils"
)

var (
	Ova = cobra.Command{
		Use:   "ova",
		Short: "Upload an image to a vSphere server",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {

			var err error
			image, _ := cmd.Flags().GetString("image")
			network, _ := cmd.Flags().GetString("network")
			name, _ := cmd.Flags().GetString("name")
			if name == "" {
				name = utils.GetBaseName(image)
			}

			ext := path.Ext(image)

			if ext != ".ova" {
				image, err = Create(name, image, make(map[string]string))
				if err != nil {
					log.Fatalf("Failed to create OVA %s", err)
				}
			}

			if err := Import(name, image, network); err != nil {
				log.Fatalf("Failed to upload %s: %v", name, err)
			}
		},
	}
)

func init() {
	Ova.Flags().String("image", "", "A local or remote path to a disk image")
	Ova.Flags().String("name", "", "Name of the template")
	Ova.Flags().String("folder", utils.GetEnvOrDefault("GOVC_FOLDER", "vm"), " The folder to upload to")
	Ova.Flags().String("network", utils.GetEnvOrDefault("GOVC_NETWORK", "VM Network"), " The VM network")
}
