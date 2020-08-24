package phases

import (
	"fmt"
	"path"

	. "github.com/flanksource/konfigadm/pkg/types" // nolint: golint
)

var AnsiblePhase Phase = ansible{}

type ansible struct{}

func (p ansible) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	var commands []Command
	files := Filesystem{}

	if len(sys.Ansible) > 0 {
		sys.AppendPackages(nil, Package{
			Name: "ansible",
		})
		sys.AppendPackages(nil, Package{
			Name:  "python",
			Flags: []Flag{DEBIAN},
		})
		sys.AppendPackageRepo(PackageRepo{
			Name:   "epel",
			URL:    "http://download.fedoraproject.org/pub/epel/7/\\$basearch",
			GPGKey: "https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7",
		}, REDHAT_LIKE)
	}

	for _, playbook := range sys.Ansible {
		filename := path.Join(playbook.Workspace, playbook.PlaybookPath)
		files[filename] = File{
			Content:     string(playbook.Playbook),
			Permissions: "0600",
			Owner:       "root",
		}
		sys.AddCommand(fmt.Sprintf("cd %s && ansible-playbook -i 'localhost, ' %s", playbook.Workspace, playbook.PlaybookPath))
	}
	return commands, files, nil
}

func (p ansible) Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool {
	verify := true
	return verify
}
