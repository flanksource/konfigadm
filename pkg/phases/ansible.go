package phases

import (
	"fmt"
	"path"

	"github.com/flanksource/konfigadm/pkg/types"
)

var AnsiblePhase types.Phase = ansible{}

type ansible struct{}

func (p ansible) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	var commands []types.Command
	files := types.Filesystem{}

	if len(sys.Ansible) > 0 {
		sys.AppendPackages(nil, types.Package{
			Name: "ansible",
		})
		sys.AppendPackages(nil, types.Package{
			Name:  "python",
			Flags: []types.Flag{types.DEBIAN},
		})
		sys.AppendPackageRepo(types.PackageRepo{
			Name:   "epel",
			URL:    "http://download.fedoraproject.org/pub/epel/7/\\$basearch",
			GPGKey: "https://dl.fedoraproject.org/pub/epel/RPM-GPG-KEY-EPEL-7",
		}, types.REDHAT_LIKE)
	}

	for _, playbook := range sys.Ansible {
		filename := path.Join(playbook.Workspace, playbook.PlaybookPath)
		files[filename] = types.File{
			Content:     string(playbook.Playbook),
			Permissions: "0600",
			Owner:       "root",
		}
		sys.AddCommand(fmt.Sprintf("cd %s && ansible-playbook -i 'localhost, ' %s", playbook.Workspace, playbook.PlaybookPath))
	}
	return commands, files, nil
}

func (p ansible) Verify(cfg *types.Config, results *types.VerifyResults, flags ...types.Flag) bool {
	verify := true
	return verify
}
