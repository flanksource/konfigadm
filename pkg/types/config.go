package types

import (
	"fmt"
	"io/ioutil"
	goos "os"
	"reflect"
	"strings"

	"github.com/moshloop/configadm/pkg/os"

	cloudinit "github.com/moshloop/configadm/pkg/cloud-init"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	yaml "gopkg.in/yaml.v3"
)

var (
	Dig = dig.New()
)

func (sys *Config) Verify(results *VerifyResults) bool {
	var Phases *[]Phase
	err := Dig.Invoke(func(_phases *[]Phase) {
		Phases = _phases
	})
	if err != nil {
		log.Fatal(err)
	}
	verify := true
	for _, phase := range *Phases {
		switch v := phase.(type) {
		case VerifyPhase:
			log.Tracef("Verifying %s", reflect.TypeOf(phase).Name())
			_verify := v.Verify(sys, results, sys.Context.Flags...)
			log.Tracef("%s -> %s", reflect.TypeOf(phase).Name(), _verify)
			verify = verify && _verify
		}

	}
	return verify
}

func (sys *Config) ApplyPhases() (Filesystem, []Command, error) {
	var Phases *[]Phase
	err := Dig.Invoke(func(_phases *[]Phase) {
		Phases = _phases
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, phase := range *Phases {
		log.Tracef("Processing flags %s(%s)", reflect.TypeOf(phase).Name(), sys.Context.Flags)
		switch v := phase.(type) {
		case ProcessFlagsPhase:
			v.ProcessFlags(sys, sys.Context.Flags...)
		}

	}

	files := Filesystem{}
	commands := sys.PreCommands

	for _, phase := range *Phases {
		c, f, err := phase.ApplyPhase(sys, sys.Context)
		log.Tracef("Applied phase %s: %s/%s", reflect.TypeOf(phase).Name(), c, f)

		if err != nil {
			return nil, []Command{}, err
		}
		for k, v := range f {
			files[k] = v
		}
		commands = append(commands, c...)
	}
	commands = append(commands, sys.Commands...)
	commands = append(commands, sys.PostCommands...)

	log.Tracef("Commands before filtering %+v\n", commands)

	//Apply tag filters on any output commands
	commands = FilterFlags(commands, sys.Context.Flags...)
	log.Tracef("Commands after filtering %+v\n", commands)

	return files, commands, nil
}

//ToCloudInit will apply all phases and produce a CloudInit object from the results
func (sys *Config) ToCloudInit() cloudinit.CloudInit {
	cloud := sys.Extra

	files, commands, err := sys.ApplyPhases()
	if err != nil {
		log.Fatal(err)
	}

	for path, content := range files {
		cloud.AddFile(path, content.Content)
	}
	cloud.AddFile(fmt.Sprintf("/usr/bin/%s.sh", Configadm), ToScript(commands))
	cloud.AddCommand(fmt.Sprintf("/usr/bin/%s.sh", Configadm))
	return *cloud
}

//ToScript returns a bash script of all the commands that can be run directly
func ToScript(commands []Command) string {
	script := "#!/bin/bash\n"
	for _, command := range commands {
		script += command.Cmd + "\n"
	}
	return script
}

func (sys *Config) Init() {
	sys.Services = make(map[string]Service)
	sys.Extra = &cloudinit.CloudInit{}
	sys.Environment = make(map[string]string)
	sys.Files = make(map[string]string)
	sys.Templates = make(map[string]string)
	sys.Sysctls = make(map[string]string)
	sys.PackageRepos = &[]PackageRepo{}
	sys.Packages = &[]Package{}
	sys.Context = &SystemContext{
		Name: Configadm,
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

func (builder *ConfigBuilder) Build() (*Config, error) {
	cfg := &Config{}
	cfg.Init()
	cfg.Context.Flags = builder.flags
	for _, config := range builder.configs {
		c, err := newConfig(config)
		if err != nil {
			log.Fatalf("Error parsing %s: %s", config, err)
		}
		cfg.ImportConfig(*c)
	}

	for _, _os := range os.SupportedOperatingSystems {
		if _os.DetectAtRuntime() && cfg.Context.OS == nil {
			cfg.Context.OS = _os
			break
		}

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

func newConfig(config string) (*Config, error) {
	c := &Config{}
	c.Init()
	if config == "-" {
		data, _ := ioutil.ReadAll(goos.Stdin)
		if err := yaml.Unmarshal(data, &c); err != nil {
			return nil, fmt.Errorf("Error reading from stdin: %s", err)
		}
		return c, nil
	}
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

func (sys Config) String() {
}

//ImportConfig merges to configs together, everything but containerRuntime and Kubernetes configs are merged
func (sys *Config) ImportConfig(c2 Config) {
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
	pkgRepos := append(*sys.PackageRepos, *c2.PackageRepos...)
	sys.PackageRepos = &pkgRepos
	pkgs := append(*sys.Packages, *c2.Packages...)
	sys.Packages = &pkgs
	sys.Timezone = c2.Timezone
	sys.ContainerRuntime = c2.ContainerRuntime
	sys.Kubernetes = c2.Kubernetes
}
