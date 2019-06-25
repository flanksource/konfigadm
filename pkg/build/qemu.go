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

	config.PostCommands = append(config.PostCommands, types.Command{Cmd: "cloud-init clean; shutdown -h now"})
	iso, _ := cloudinit.CreateISO("builder", config.ToCloudInit().String())
	cmdLine := fmt.Sprintf(`
	kvm -M pc -m 1024 -smp 2 \
	-nographic \
	-hda %s \
	-drive "file=%s,if=virtio,format=raw" \
	-net nic -net user,hostfwd=tcp:127.0.0.1:2022-:22`, image, iso)

	log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
	if err := utils.Exec(cmdLine); err != nil {
		log.Fatalf("Failed to run: %s, %s", cmdLine, err)
	}
}
