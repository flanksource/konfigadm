package phases

import (
	"io/ioutil"
	"os"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
	log "github.com/sirupsen/logrus"
)

var Files Phase = filesPhase{}

type filesPhase struct{}

func (p filesPhase) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := make(map[string]File)
	for k, v := range sys.Filesystem {
		files[k] = v
	}

	for k, v := range sys.Files {
		files[k] = File{Content: Lookup(v)}
	}

	return commands, files, nil
}

func (p filesPhase) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	for f := range cfg.Files {

		if _, err := os.Stat(f); err != nil {
			verify = false
			results.Fail("%s does not exist", f)
		} else {
			results.Pass("%s exists", f)
		}
	}

	for f := range cfg.Templates {
		if _, err := os.Stat(f); err != nil {
			verify = false
			results.Fail("%s does not exist", f)
		} else {
			results.Pass("%s exists", f)
		}
	}

	return verify
}

func Lookup(path string) string {
	_, err := os.Stat(path)
	if err != nil {
		return path
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading: %s : %s", path, err)
		return path
	}
	return string(data)
}
