package build

import (
	"fmt"
	"os"

	"github.com/mitchellh/colorstring"
	"github.com/moshloop/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Libguestfs struct{}

func (l Libguestfs) Build(image string, config *os.File) {
	os.Remove("/tmp/builder.log")
	os.Remove("builder.log")

	konfigadmPath := os.Args[0]
	cmdLine := fmt.Sprintf("virt-customize -a %s --delete /tmp/builder.log --copy-in %s:/tmp --copy-in %s:/tmp --run-command '/tmp/konfigadm apply -vv -c %s' ", image, konfigadmPath, config.Name(), config.Name())

	log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
	if err := utils.Exec(cmdLine); err != nil {
		log.Errorf("builder.log: %s\n", utils.SafeRead("builder.log"))
		log.Fatalf("Failed to run: %s, %s", cmdLine, err)
	}
	utils.Exec(fmt.Sprintf("virt-copy-out -a %s /tmp/builder.log .", image))
	log.Infof("builder.log %s\n", utils.SafeRead("builder.log"))
}
