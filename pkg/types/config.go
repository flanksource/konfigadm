package types

import (
	"fmt"
	"io/ioutil"
	goos "os"
	"reflect"
	"strings"

	cloudinit "github.com/moshloop/konfigadm/pkg/cloud-init"
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
			verify = verify && v.Verify(sys, results, sys.Context.Flags...)
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
		log.Tracef("Applied phase %s: %s/%v", reflect.TypeOf(phase).Name(), c, f)

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

	log.Tracef("Files before filtering: %s\n", GetKeys(files))
	files = FilterFilesystemByFlags(files, sys.Context.Flags...)
	log.Tracef("Files after filtering: %s\n", GetKeys(files))
	return files, commands, nil
}

func GetKeys(m map[string]File) []string {
	keys := []string{}
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

//ToCloudInit will apply all phases and produce a CloudInit object from the results
func (sys *Config) ToCloudInit() cloudinit.CloudInit {
	cloud := sys.Extra
	log.Debugf("Extra: %+v", cloud)

	files, commands, err := sys.ApplyPhases()
	if err != nil {
		log.Fatal(err)
	}

	for path, content := range files {
		cloud.AddFile(path, content.Content)
	}
	yml, _ := yaml.Marshal(sys)
	cloud.AddFile("/etc/konfigadm.yml", string(yml))
	cloud.AddFile(fmt.Sprintf("/usr/bin/%s.sh", konfigadm), ToScript(commands))
	cloud.AddCommand(fmt.Sprintf("/usr/bin/%s.sh", konfigadm))
	return *cloud
}

//ToScript returns a bash script of all the commands that can be run directly
func ToScript(commands []Command) string {
	script := "#!/bin/bash\nset -o verbose\n"
	for _, command := range commands {
		script += command.Cmd + "\n"
	}
	return script
}

func (sys *Config) Init() {
	sys.Services = make(map[string]Service)
	if sys.Extra == nil {
		sys.Extra = &cloudinit.CloudInit{}
	}
	sys.Environment = make(map[string]string)
	sys.Files = make(map[string]string)
	sys.Filesystem = make(map[string]File)
	sys.Templates = make(map[string]string)
	sys.Sysctls = make(map[string]string)
	sys.PackageRepos = &[]PackageRepo{}
	sys.Packages = &[]Package{}
	sys.Context = &SystemContext{
		Name: konfigadm,
		Vars: make(map[string]interface{}),
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
	if c2.Extra != nil {
		if strings.TrimSpace(fmt.Sprintf("%+v", sys.Extra)) == "{}" {
			sys.Extra = c2.Extra
		} else if strings.TrimSpace(fmt.Sprintf("%+v", c2.Extra)) != "#cloud-config\n{}" {
			log.Warnf("More than 1 extra cloud-init section found , merging clout-init is not supported and will override")
			sys.Extra = c2.Extra
		}
	}
	if c2.Cleanup != nil {
		sys.Cleanup = c2.Cleanup
	}
	sys.Commands = append(sys.Commands, c2.Commands...)
	sys.PreCommands = append(sys.PreCommands, c2.PreCommands...)
	sys.PostCommands = append(sys.PostCommands, c2.PostCommands...)
	sys.Users = append(sys.Users, c2.Users...)

	for k, v := range c2.Files {
		sys.Files[k] = v
	}
	for k, v := range c2.Filesystem {
		sys.Filesystem[k] = v
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
	if c2.ContainerRuntime != nil {
		sys.ContainerRuntime = c2.ContainerRuntime
	}
	if c2.Kubernetes != nil {
		sys.Kubernetes = c2.Kubernetes
	}
}
