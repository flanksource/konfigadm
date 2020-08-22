package build

import (
	"io/ioutil"
	"os"
	"runtime"

	"github.com/flanksource/konfigadm/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Scratch interface {
	Create() error
	UnwrapToDir(dir string) error
	GetImg() string
}

type DarwinScratch struct {
	img string
}

func NewScratch() Scratch {
	var scratch Scratch
	if runtime.GOOS == "darwin" {
		scratch = &DarwinScratch{}
	}
	if err := scratch.Create(); err != nil {
		log.Errorf("Failed to create: %s", err)
	}
	return scratch
}
func (s *DarwinScratch) GetImg() string {
	return s.img
}
func (s *DarwinScratch) Create() error {
	tmp, _ := ioutil.TempFile("", "scratch*.img")
	s.img = tmp.Name()
	log.Infof("Creating %s", s.img)
	if err := utils.Exec("hdiutil create -fs FAT32 -size 100m    -volname scratch %s", s.img); err != nil {
		return err
	}
	return os.Rename(s.img+".dmg", s.img)
}

func (s *DarwinScratch) UnwrapToDir(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	mount, _ := ioutil.TempDir("", "mount")
	if err := utils.Exec("hdiutil attach -mountpoint %s  %s ", mount, s.img); err != nil {
		return err
	}
	defer utils.Exec("hdiutil detach %s", mount) // nolint: errcheck
	return utils.Exec("cp -r %s/* %s", mount, dir)
}

func CaptureLogCommands() []string {
	return []string{
		"mkdir -p /scratch",
		"journalctl  -b --no-hostname -o short > /scratch/journal.log",
		"cp -r /var/log/ /scratch || true",
		"cp -r /var/lib/cloud/ /scratch || true",
	}
}
