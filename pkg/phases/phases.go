package phases

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	cloudinit "github.com/moshloop/configadm/pkg/cloud-init"
	. "github.com/moshloop/configadm/pkg/utils"
	yaml "gopkg.in/yaml.v3"
)

func (sys *SystemConfig) Init() {
	sys.Services = make(map[string]Service)
	sys.Extra = &cloudinit.CloudInit{}
	sys.Environment = make(map[string]string)
	sys.Files = make(map[string]string)
	sys.Templates = make(map[string]string)
	sys.Sysctls = make(map[string]string)
	sys.Packages = Packages{}
	sys.Context = &SystemContext{
		Name: "cloud-config",
		Vars: make(map[string]interface{}),
	}
}

type ConfigBuilder struct {
	configs []string
	vars    []string
	flags   []Flag
}

func (f *ConfigBuilder) WithVars(vars ...string) *ConfigBuilder {
	f.vars = vars
	return f
}

func (f *ConfigBuilder) WithFlags(flags ...Flag) *ConfigBuilder {
	f.flags = flags
	return f
}

func (builder *ConfigBuilder) Build() (*SystemConfig, error) {
	cfg := &SystemConfig{}
	cfg.Init()
	for _, config := range builder.configs {
		c, err := newSystemConfig(config)
		if err != nil {
			log.Fatalf("Error parsing %s: %s", config, err)
		}
		cfg.ImportConfig(*c)
	}

	for _, v := range builder.vars {
		if strings.Contains(v, "=") {
			cfg.Context.Vars[strings.Split(v, "=")[0]] = strings.Split(v, "=")[1]
		}
	}

	if len(builder.flags) > 0 {
		cfg.MinifyWithFlags(builder.flags...)
	}
	cfg.Transform(*cfg.Context)
	return cfg, nil
}

func NewConfig(configs ...string) *ConfigBuilder {
	return &ConfigBuilder{
		configs: configs,
	}
}

func newSystemConfig(config string) (*SystemConfig, error) {
	c := &SystemConfig{}
	c.Init()
	data, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(config, "yml") {
		err = yaml.Unmarshal(data, &c)
	} else {
		return nil, fmt.Errorf("Unknown file type: %s", config)
	}

	return c, err
}

func (sys *SystemConfig) ToFiles() map[string]string {
	files := make(map[string]string)
	for k, v := range sys.Files {
		files[k] = v
	}

	if len(sys.Services) > 0 {
		for name, svc := range sys.Services {
			filename := fmt.Sprintf("/etc/systemd/system/%s.service", name)
			files[filename] = svc.Extra.ToUnitFile()
			sys.Commands = append(sys.Commands, Command{Cmd: "systemctl enable " + name})
			sys.Commands = append(sys.Commands, Command{Cmd: "systemctl start " + name})
		}
	}

	if len(sys.Environment) > 0 {
		files["/etc/environment"] = MapToIni(sys.Environment)
	}

	return files
}

func (sys SystemConfig) ToScript() string {
	script := ""
	for _, cmd := range sys.PreCommands {
		script += cmd.Cmd + "\n"
	}
	for k, v := range sys.Environment {
		script += fmt.Sprintf("export %s=\"%s\"\n", k, v)
	}
	for _, cmd := range sys.Commands {
		script += cmd.Cmd + "\n"
	}
	for _, cmd := range sys.PostCommands {
		script += cmd.Cmd + "\n"
	}
	return script
}

func (sys SystemConfig) ToCloudInit() cloudinit.CloudInit {
	cloud := sys.Extra

	for path, content := range sys.ToFiles() {
		cloud.AddFile(path, content)
	}
	cloud.AddFile("/usr/bin/cloud-config.sh", sys.ToScript())
	cloud.AddCommand("/usr/bin/cloud-config.sh")
	return *cloud
}

func (sys SystemConfig) String() {
}

//ImportConfig merges to configs together, everything but containerRuntime and Kubernetes configs are merged
func (cfg *SystemConfig) ImportConfig(c2 SystemConfig) {
	cfg.Commands = append(cfg.Commands, c2.Commands...)
	cfg.PreCommands = append(cfg.PreCommands, c2.PreCommands...)
	cfg.PostCommands = append(cfg.PostCommands, c2.PostCommands...)
	cfg.Users = append(cfg.Users, c2.Users...)

	for k, v := range c2.Files {
		cfg.Files[k] = v
	}
	for k, v := range c2.Templates {
		cfg.Templates[k] = v
	}
	for k, v := range c2.Environment {
		cfg.Environment[k] = v
	}
	for k, v := range c2.Services {
		cfg.Services[k] = v
	}
	for k, v := range c2.Sysctls {
		cfg.Sysctls[k] = v
	}

	cfg.Containers = append(cfg.Containers, c2.Containers...)
	cfg.Images = append(cfg.Images, c2.Images...)
	cfg.PackageRepos = append(cfg.PackageRepos, c2.PackageRepos...)
	cfg.Packages = append(cfg.Packages, c2.Packages...)
	cfg.Timezone = c2.Timezone
	cfg.ContainerRuntime = c2.ContainerRuntime
	cfg.Kubernetes = c2.Kubernetes
}
