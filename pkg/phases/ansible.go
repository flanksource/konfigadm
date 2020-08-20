package phases

import (
	"fmt"
	"path"

	. "github.com/flanksource/konfigadm/pkg/types"
)

var Ansible Phase = ansible{}

type ansible struct{}

func (p ansible) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	if len(sys.AnsiblePlaybooks) > 0 {
		sys.AddPackage("ansible", nil)
		sys.AddPackage("python", &DEBIAN)
		sys.AppendPackageRepo(PackageRepo{
			Name:   "docker-ce",
			URL:    "http://download.fedoraproject.org/pub/epel/7/\\$basearch",
			GPGKey: "https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7",
		}, REDHAT_LIKE)
	}

	for i, playbook := range sys.AnsiblePlaybooks {
		playbookName := fmt.Sprintf("playbook-%d.yml", i)
		filename := path.Join("/tmp", "ansible-playbooks", playbookName)
		files[filename] = File{
			Content:     string(playbook),
			Permissions: "0600",
			Owner:       "root",
		}
		commands = append(commands, Command{Cmd: fmt.Sprintf("ansible-playbook -i 'localhost, ' %s", filename)})
		sys.AddCommand(fmt.Sprintf("ansible-playbook -i 'localhost, ' %s", filename))
	}
	return commands, files, nil
}

func (p ansible) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	return verify
}
