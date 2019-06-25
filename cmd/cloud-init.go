package cmd

import (
	"fmt"
	"log"

	cloudinit "github.com/moshloop/konfigadm/pkg/cloud-init"
	"github.com/spf13/cobra"
)

var (
	CloudInit = cobra.Command{
		Use:   "cloud-init",
		Short: "Exports the configuration in cloud-init format",
		Args:  cobra.MinimumNArgs(0),

		Run: func(cmd *cobra.Command, args []string) {

			cfg := GetConfig(cmd, args)
			var userdata string
			if base64, _ := cmd.Flags().GetBool("base64"); base64 {
				cfg.Extra.FileEncoding = "base64"
				userdata = cfg.ToCloudInit().String()
			} else {
				userdata = cfg.ToCloudInit().String()
			}

			if iso, _ := cmd.Flags().GetBool("iso"); iso {
				hostname, _ := cmd.Flags().GetString("hostname")
				path, err := cloudinit.CreateISO(hostname, userdata)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(path)

			} else {
				fmt.Println(userdata)
			}
		},
	}
)

func init() {
	CloudInit.Flags().Bool("base64", true, "Base64 encode files")
	CloudInit.Flags().Bool("iso", false, "Create an ISO with the cloud-init embedded as user-metadata")
	CloudInit.Flags().String("hostname", "", "Hostname to use in generated cloud-init")
}
