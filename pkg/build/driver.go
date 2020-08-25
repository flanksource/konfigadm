package build

import (
	"fmt"
	"strings"

	cloudinit "github.com/flanksource/konfigadm/pkg/cloud-init"
	"github.com/flanksource/konfigadm/pkg/types"
)

type Driver interface {
	Build(image string, cfg *types.Config)
	Test(image string, cfg *types.Config, privateKeyFile string, template string) error
}

func createIso(config *types.Config) (string, error) {
	cloudInit := config.ToCloudInit()

	if config.Context.CaptureLogs != "" {
		cloudInit.Runcmd = append([][]string{[]string{"bash", "-x", "-c", "mkdir /scratch; mount /dev/sdb1 /scratch"}}, cloudInit.Runcmd...)
	}
	if config.Context.CaptureLogs != "" && (config.Cleanup == nil || !*config.Cleanup) {
		cloudInit.Runcmd = append(cloudInit.Runcmd, []string{"bash", "-x", "-c", strings.Join(CaptureLogCommands(), "; ")})
	}

	// PowerState is once per instance and cloud-init clean (creating a new instance) fails on ubuntu 18.04:
	// IsADirectory: /var/lib/cloud/instance
	//	"cloudInit.PowerState.Mode = "poweroff"
	// so we append a shutdown manually
	cloudInit.Runcmd = append(cloudInit.Runcmd, []string{"shutdown", "-h", "now"})
	return cloudinit.CreateISO("builder", cloudInit.String())
}

func createTestIso(config *types.Config) (string, error) {
	cloudInit := config.ToCloudInit()

	if config.Context.CaptureLogs != "" {
		cloudInit.Runcmd = append([][]string{[]string{"bash", "-x", "-c", "mkdir /scratch; mount /dev/sdb1 /scratch"}}, cloudInit.Runcmd...)
	}
	if config.Context.CaptureLogs != "" && (config.Cleanup == nil || !*config.Cleanup) {
		cloudInit.Runcmd = append(cloudInit.Runcmd, []string{"bash", "-x", "-c", strings.Join(CaptureLogCommands(), "; ")})
	}

	fmt.Printf("Cloud init: \n%s\n", cloudInit)
	return cloudinit.CreateISO("builder", cloudInit.String())
}
