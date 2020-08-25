package phases

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/flanksource/konfigadm/pkg/types"
)

var Users types.Phase = users{}

type users struct{}

func (u users) ApplyPhase(sys *types.Config, ctx *types.SystemContext) ([]types.Command, types.Filesystem, error) {
	files := types.Filesystem{}
	var commands []types.Command

	for _, user := range sys.Users {

		if user.Sudo != "" {
			files["/etc/sudoers.d/91-"+user.Name] = types.File{Content: fmt.Sprintf("%s %s", user.Name, user.Sudo)}
		}
		cmd := fmt.Sprintf("getent passwd %s || (useradd -m", user.Name)

		if user.Shell != "" {
			cmd += " -s " + user.Shell
		}
		if user.UID != "" {
			cmd += " -u " + user.UID
		}

		if user.Gecos != "" {
			cmd += fmt.Sprintf(" -c \"%s\"", user.Gecos)
		}
		cmd += fmt.Sprintf(" %s ) ", user.Name)

		authorizedKeys := base64.StdEncoding.EncodeToString([]byte(strings.Join(user.SSHAuthorizedKeys, "\n")))

		commands = append(commands, types.Command{Cmd: cmd})
		commands = append(commands, types.Command{Cmd: fmt.Sprintf("mkdir -p /home/%s/.ssh/ && ( echo %s | base64 -d > /home/%s/.ssh/authorized_keys ) && chown %s /home/%s/.ssh", user.Name, authorizedKeys, user.Name, user.Name, user.Name)})

	}
	return commands, files, nil

}
