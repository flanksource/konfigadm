package phases

import (
	"io/ioutil"
	"os"

	. "github.com/moshloop/configadm/pkg/types"
	log "github.com/sirupsen/logrus"
)

var Files Phase = filesPhase{}

type filesPhase struct{}

func (p filesPhase) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	for k, v := range sys.Files {
		files[k] = File{Content: Lookup(v)}
	}

	return commands, files, nil
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
