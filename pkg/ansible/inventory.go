package ansible

import (
	"bufio"
	"path/filepath"
	"os"
	"strings"
	log "github.com/sirupsen/logrus"
)

func ParseInventory(inventory string) []string {
	inventory = filepath.Clean(inventory)
	log.Infof("Reading inventory %s\n", inventory)
	f, err := os.Open(inventory)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = f.Close(); err != nil {
			panic(err)
		}
	}()

	s := bufio.NewScanner(f)
	remoteHosts := make([]string, 0)
	isChildGroup := false
	for s.Scan() {
		line := strings.Split(s.Text(), " ")[0] // Ignore all set host variables
		if strings.HasSuffix(line, ":children]") { // Ignore child groups completely
			isChildGroup = true
			continue
		} else if strings.HasPrefix(line, "[") { // Ignore group names
			isChildGroup = false // Ends any previous child group
			continue
		}
		if isChildGroup {
			continue
		}
		if line == "" {
			continue
		}
		duplicate := false
		for _, host := range remoteHosts {
			if host == line {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		remoteHosts = append(remoteHosts, line)
	}
	return remoteHosts
}
