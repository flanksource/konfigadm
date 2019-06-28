package build

import (
	"fmt"

	cloudinit "github.com/moshloop/konfigadm/pkg/cloud-init"
	"github.com/moshloop/konfigadm/pkg/types"

	"github.com/mitchellh/colorstring"
	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Qemu struct{}

func (q Qemu) Build(image string, config *types.Config) {
	cloud_init := config.ToCloudInit()
	cloud_init.PowerState.Mode = "poweroff"
	iso, _ := cloudinit.CreateISO("builder", cloud_init.String())
	cmdLine := fmt.Sprintf(`qemu-system-x86_64 \
				-global virtio-blk-pci.scsi=off \
				-enable-kvm \
		    -enable-fips \
		    -nodefaults \
		    -display none \
		    -machine accel=kvm \
		    -cpu host -smp cpus=2 \
		    -m 1024 \
		    -no-reboot \
		    -rtc driftfix=slew \
		    -no-hpet \
		    -global kvm-pit.lost_tick_policy=discard \
		    -object rng-random,filename=/dev/urandom,id=rng0 \
		    -device virtio-rng-pci,rng=rng0 \
				-hda %s \
		    -cdrom %s \
		    -device virtio-serial-pci \
		    -serial stdio \
				-device sga \
				-net nic -net user,hostfwd=tcp:127.0.0.1:2022-:22`, image, iso)

	log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
	if err := utils.Exec(cmdLine); err != nil {
		log.Fatalf("Failed to run: %s, %s", cmdLine, err)
	}
}
