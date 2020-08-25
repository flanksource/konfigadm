package phases

import (
	"github.com/flanksource/konfigadm/pkg/types"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

var Files types.Phase = filesPhase{}

type filesPhase struct{}

func (p filesPhase) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	var commands []types.Command
	files := make(map[string]types.File)
	for k, v := range sys.Filesystem {
		files[k] = v
	}

	for k, v := range sys.Files {
		files[k] = types.File{Content: Lookup(v)}
	}

	return commands, files, nil
}

func (p filesPhase) Verify(cfg *types.Config, results *types.VerifyResults, flags ...types.Flag) bool {
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
