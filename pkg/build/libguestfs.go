package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
	"github.com/mitchellh/colorstring"
	log "github.com/sirupsen/logrus"
	"gopkg.in/flanksource/yaml.v3"
)

type Libguestfs struct{}

func (l Libguestfs) Build(image string, cfg *types.Config) {

	_, _, err := cfg.ApplyPhases()
	if err != nil {
		log.Fatalf("Error applying phases %s\n", err)
	}
	data, _ := yaml.Marshal(cfg)
	tmpfile, err := ioutil.TempFile("", "konfigadm.*.yml")
	if err != nil {
		log.Fatalf("Cannot create tempfile %s", err)
	}
	if _, err := tmpfile.Write(data); err != nil {
		log.Fatalf("Error writing tmp file %s", err)
	}
	os.Remove("/tmp/builder.log")
	os.Remove("builder.log")

	konfigadmPath, _ := os.Executable()
	konfigAdm := path.Base(konfigadmPath)
	cmdLine := fmt.Sprintf(`
		virt-customize -a %s \
		--delete /tmp/builder.log \
		--copy-in %s:/tmp \
		--copy-in %s:/tmp \
		--run-command '/tmp/%s apply -vv -c %s' 1>&2`, image, konfigadmPath, tmpfile.Name(), konfigAdm, tmpfile.Name())

	log.Infof("Executing %s\n", colorstring.Color("[light_green]"+cmdLine))
	if err := utils.Exec(cmdLine); err != nil {
		log.Errorf("builder.log: %s\n", utils.SafeRead("builder.log"))
		log.Fatalf("Failed to run: %s, %s", cmdLine, err)
	}
	utils.Exec(fmt.Sprintf("virt-copy-out -a %s /tmp/builder.log .", image))
	log.Infof("builder.log %s\n", utils.SafeRead("builder.log"))
}
