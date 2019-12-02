package build

import (
	"fmt"

	"github.com/mitchellh/colorstring"
	log "github.com/sirupsen/logrus"

	"github.com/moshloop/konfigadm/pkg/types"
	"github.com/moshloop/konfigadm/pkg/utils"
)

type Qemu struct{}

func (q Qemu) Build(image string, config *types.Config) {
	var scratch Scratch
	if config.Context.CaptureLogs != "" {
		log.Infof("Using scratch directory / disk")
		scratch = NewScratch()
	}

	iso, err := createIso(config)
	if err != nil {
		log.Fatalf("Failed to build ISO %v", err)
	}
	if iso == "" {
		log.Fatalf("Empty ISO created")
	}
	cmdLine := qemuSystem(image, iso)
	if config.Context.CaptureLogs != "" {
		cmdLine += fmt.Sprintf(" -hdb %s", scratch.GetImg())
	}

	log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
	if err := utils.Exec(cmdLine); err != nil {
		log.Fatalf("Failed to run: %s, %s", cmdLine, err)
	}
	if config.Context.CaptureLogs != "" {
		log.Infof("Coping captured logs to %s\n", config.Context.CaptureLogs)
		scratch.UnwrapToDir(config.Context.CaptureLogs)
	}
}

func qemuSystem(disk, iso string) string {
	return fmt.Sprintf(`qemu-system-x86_64 \
		-nodefaults \
		-display none \
		-machine accel=kvm:hvf \
		-cpu host -smp cpus=2 \
		-m 1024 \
		-hda %s \
		-cdrom %s \
		-device virtio-serial-pci \
		-serial stdio \
		-net nic -net user,hostfwd=tcp:127.0.0.1:2022-:22`, disk, iso)
}
