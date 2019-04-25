package cloudmeta

import (
	"io/ioutil"
	"log"
	"os"
)

func init() {
	Register(TransformFiles)
}

func TransformFiles(sys *SystemConfig, ctx *SystemContext) (commands []string, files map[string]string, err error) {
	commands = []string{}
	files = make(map[string]string)

	for k, v := range sys.Files {
		files[k] = Lookup(v)
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
