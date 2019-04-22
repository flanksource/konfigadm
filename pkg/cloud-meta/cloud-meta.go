package cloudmeta

import (
	"fmt"
	"strings"

	. "github.com/moshloop/cloud-config/pkg/cloud-init"
	. "github.com/moshloop/cloud-config/pkg/utils"
)

func (sys *SystemConfig) Init() {
	sys.Services = make(map[string]Service)
	sys.Extra = CloudInit{}
}

func (sys *SystemConfig) ToFiles() map[string]string {
	files := make(map[string]string)
	for k, v := range sys.Files {
		files[k] = v
	}

	if len(sys.Sysctls) > 0 {
		filename := fmt.Sprintf("/etc/sysctl.conf.d/100-%s.conf", sys.Context.Name)
		files[filename] = MapToIni(sys.Sysctls)
	}

	if len(sys.Services) > 0 {
		for name, svc := range sys.Services {
			filename := fmt.Sprintf("/etc/systemd/system/%s.service", name)
			files[filename] = svc.Extra.ToUnitFile()
			sys.Commands = append(sys.Commands, "systemctl enable "+name)
			sys.Commands = append(sys.Commands, "systemctl start "+name)

		}
	}

	if len(sys.Environment) > 0 {
		files["/etc/environment"] = MapToIni(sys.Environment)
	}

	return files
}

func (sys SystemConfig) ToScript() string {
	script := ""
	script += strings.Join(sys.PreCommands, "\n")
	script += strings.Join(sys.Commands, "\n")
	script += strings.Join(sys.PostCommands, "\n")
	return script
}

func (sys SystemConfig) ToCloudInit() CloudInit {
	cloud := sys.Extra

	for path, content := range sys.Files {
		cloud.AddFile(path, content)
	}
	for path, content := range sys.ToFiles() {
		cloud.AddFile(path, content)
	}
	cloud.AddFile("/usr/bin/cloudinit.sh", sys.ToScript())
	cloud.AddCommand("/usr/bin/cloudinit.sh")
	return cloud
}

func (sys SystemConfig) String() {

}
