package phases

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	cloudinit "github.com/moshloop/configadm/pkg/cloud-init"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

var (
	Phases = []Phase{
		Context,
		Sysctl,
		Environment,
		Containers,
		Packages,
		Services,
		Files,
		Commands,
	}
)

func (sys *SystemConfig) ApplyPhases() (files Filesystem, script string, err error) {
	for _, phase := range Phases {
		log.Tracef("Processing flags %s(%s)", reflect.TypeOf(phase).Name(), sys.Context.Flags)
		switch v := phase.(type) {
		case AllPhases:
			v.ProcessFlags(sys, sys.Context.Flags...)
		}

	}

	files = Filesystem{}
	commands := []Command{}

	for _, phase := range Phases {
		c, f, err := phase.ApplyPhase(sys, sys.Context)
		log.Tracef("Applied phase %s: %s/%s", reflect.TypeOf(phase).Name(), c, f)

		if err != nil {
			return nil, "", err
		}
		for k, v := range f {
			files[k] = v
		}
		commands = append(commands, c...)
	}

	//Apply flag filters on any output commands
	commands = filter(commands, sys.Context.Flags...)

	return files, sys.toScript(commands...), nil
}

func (sys *SystemConfig) Init() {
	sys.Services = make(map[string]Service)
	sys.Extra = &cloudinit.CloudInit{}
	sys.Environment = make(map[string]string)
	sys.Files = make(map[string]string)
	sys.Templates = make(map[string]string)
	sys.Sysctls = make(map[string]string)
	sys.Packages = []Package{}
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
	cfg.Context.Flags = builder.flags
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

func (sys SystemConfig) toScript(commands ...Command) string {
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
	for _, cmd := range commands {
		script += cmd.Cmd + "\n"
	}
	for _, cmd := range sys.PostCommands {
		script += cmd.Cmd + "\n"
	}
	return script
}

//ToCloudInit will apply all phases and produce a CloudInit object from the results
func (sys SystemConfig) ToCloudInit() cloudinit.CloudInit {
	cloud := sys.Extra

	files, script, err := sys.ApplyPhases()
	if err != nil {
		log.Fatal(err)
	}

	for path, content := range files {
		cloud.AddFile(path, content.Content)
	}
	cloud.AddFile("/usr/bin/cloud-config.sh", script)
	cloud.AddCommand("/usr/bin/cloud-config.sh")
	return *cloud
}

func (sys SystemConfig) String() {
}

//ImportConfig merges to configs together, everything but containerRuntime and Kubernetes configs are merged
func (sys *SystemConfig) ImportConfig(c2 SystemConfig) {
	sys.Commands = append(sys.Commands, c2.Commands...)
	sys.PreCommands = append(sys.PreCommands, c2.PreCommands...)
	sys.PostCommands = append(sys.PostCommands, c2.PostCommands...)
	sys.Users = append(sys.Users, c2.Users...)

	for k, v := range c2.Files {
		sys.Files[k] = v
	}
	for k, v := range c2.Templates {
		sys.Templates[k] = v
	}
	for k, v := range c2.Environment {
		sys.Environment[k] = v
	}
	for k, v := range c2.Services {
		sys.Services[k] = v
	}
	for k, v := range c2.Sysctls {
		sys.Sysctls[k] = v
	}

	sys.Containers = append(sys.Containers, c2.Containers...)
	sys.Images = append(sys.Images, c2.Images...)
	sys.PackageRepos = append(sys.PackageRepos, c2.PackageRepos...)
	sys.Packages = append(sys.Packages, c2.Packages...)
	sys.Timezone = c2.Timezone
	sys.ContainerRuntime = c2.ContainerRuntime
	sys.Kubernetes = c2.Kubernetes
}
