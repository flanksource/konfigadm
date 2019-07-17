package build

import (
	"strings"

	cloudinit "github.com/moshloop/konfigadm/pkg/cloud-init"
	"github.com/moshloop/konfigadm/pkg/types"
)

type Driver interface {
	Build(image string, cfg *types.Config)
}

func createIso(config *types.Config) string {
	cloud_init := config.ToCloudInit()

	if config.Context.CaptureLogs != "" {
		cloud_init.Runcmd = append([][]string{[]string{"bash", "-x", "-c", "mkdir /scratch; mount /dev/sdb1 /scratch"}}, cloud_init.Runcmd...)
	}
	if config.Context.CaptureLogs != "" && config.Cleanup == nil || !*config.Cleanup {
		cloud_init.Runcmd = append(cloud_init.Runcmd, []string{"bash", "-x", "-c", strings.Join(CaptureLogCommands(), "; ")})
	}

	// PowerState is once per instance and cloud-init clean (creating a new instance) fails on ubuntu 18.04:
	// IsADirectory: /var/lib/cloud/instance
	//	"cloud_init.PowerState.Mode = "poweroff"
	// so we append a shutdown manually
	cloud_init.Runcmd = append(cloud_init.Runcmd, []string{"shutdown", "-h", "now"})
	iso, _ := cloudinit.CreateISO("builder", cloud_init.String())
	return iso
}
