package cloudmeta

import (
	"fmt"
	"strings"

	"github.com/moshloop/configadm/pkg/systemd"
)

func init() {
	Register(TransformContainers)
}

func TransformContainers(sys *SystemConfig, ctx *SystemContext) (commands []string, files map[string]string, err error) {
	commands = []string{}
	files = make(map[string]string)
	for _, c := range sys.Containers {
		filename := fmt.Sprintf("/etc/systemd/system/%s.service", c.Name())
		files[filename] = c.ToSystemDUnit()
		if len(c.Env) > 0 {
			files["/etc/environment."+c.Name()] = c.ToEnvironmentFile()
		}
		commands = append(commands, "systemctl enable "+c.Name())
		commands = append(commands, "systemctl start "+c.Name())
	}
	return commands, files, nil
}

func (c Container) Name() string {
	if c.Service != "" {
		return c.Service
	}
	name := strings.Split(c.Image, ":")[0]
	if strings.Contains(name, "/") {
		name = name[strings.LastIndex(name, "/")+1:]
	}
	return name
}

func (c Container) ToEnvironmentFile() string {
	s := ""
	for k, v := range c.Env {
		s += fmt.Sprintf("%s=%s\n", k, v)
	}
	return s
}

func (c Container) ToSystemDUnit() string {
	svc := systemd.DefaultSystemdService(c.Name())

	args := ""
	args += c.DockerOpts
	args += fmt.Sprintf(" --env-file /etc/environment.%s", c.Name())
	if c.Network != "" {
		args += " --network " + c.Network
	}

	for _, v := range c.Volumes {
		args += fmt.Sprintf(" -v %s", v)
	}

	for _, p := range c.Ports {
		args += fmt.Sprintf(" -p %d:%d", p.Port, p.Target)
	}

	svc.Service.ExecStart = fmt.Sprintf("/bin/docker run --rm --name %s %s %s %s", c.Name(), args, c.Image, c.Args)
	return svc.ToUnitFile()
}
