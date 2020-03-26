package types

import (
	"fmt"
	"io/ioutil"
	goos "os"
	"reflect"
	"strings"

	cloudinit "github.com/flanksource/konfigadm/pkg/cloud-init"
	log "github.com/sirupsen/logrus"
	"go.uber.org/dig"
	yaml "gopkg.in/yaml.v3"
)

var (
	Dig = dig.New()
)

func (cfg *Config) Verify(results *VerifyResults) bool {
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
			// run verification always, even if previous verifications have failed
			_verify := v.Verify(cfg, results, cfg.Context.Flags...)
			verify = verify && _verify
		}
	}
	return verify
}

func (cfg *Config) ApplyPhases() (Filesystem, []Command, error) {
	var Phases *[]Phase
	err := Dig.Invoke(func(_phases *[]Phase) {
		Phases = _phases
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, phase := range *Phases {
		log.Tracef("Processing flags %s(%s)", reflect.TypeOf(phase).Name(), cfg.Context.Flags)
		switch v := phase.(type) {
		case ProcessFlagsPhase:
			v.ProcessFlags(cfg, cfg.Context.Flags...)
		}
	}

	files := Filesystem{}
	commands := cfg.PreCommands

	for _, phase := range *Phases {
		c, f, err := phase.ApplyPhase(cfg, cfg.Context)
		log.Tracef("Applied phase %s: %s/%v", reflect.TypeOf(phase).Name(), c, f)

		if err != nil {
			return nil, []Command{}, err
		}
		for k, v := range f {
			files[k] = v
		}
		commands = append(commands, c...)
	}
	commands = append(commands, cfg.Commands...)
	commands = append(commands, cfg.PostCommands...)

	log.Tracef("Commands before filtering %+v\n", commands)
	//Apply tag filters on any output commands
	commands = FilterFlags(commands, cfg.Context.Flags...)
	log.Tracef("Commands after filtering %+v\n", commands)

	log.Tracef("Files before filtering: %s\n", GetKeys(files))
	files = FilterFilesystemByFlags(files, cfg.Context.Flags...)
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
func (cfg *Config) ToCloudInit() cloudinit.CloudInit {
	cloud := cfg.Extra
	log.Tracef("Extra: %+v", cloud)

	files, commands, err := cfg.ApplyPhases()
	if err != nil {
		log.Fatal(err)
	}

	for path, content := range files {
		cloud.AddFile(path, content.Content)
	}
	yml, _ := yaml.Marshal(cfg)
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

func (cfg *Config) Init() {
	cfg.Services = make(map[string]Service)
	if cfg.Extra == nil {
		cfg.Extra = &cloudinit.CloudInit{}
	}
	cfg.Environment = make(map[string]string)
	cfg.Files = make(map[string]string)
	cfg.Filesystem = make(map[string]File)
	cfg.Templates = make(map[string]string)
	cfg.Sysctls = make(map[string]string)
	cfg.PackageRepos = &[]PackageRepo{}
	cfg.Packages = &[]Package{}
	cfg.Context = &SystemContext{
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
			return nil, fmt.Errorf("error reading from stdin: %s", err)
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
		return nil, fmt.Errorf("unknown file type: %s", config)
	}

	return c, err
}

func (cfg Config) String() {
}

//ImportConfig merges to configs together, everything but containerRuntime and Kubernetes configs are merged
func (cfg *Config) ImportConfig(c2 Config) {
	if c2.Extra != nil {
		if strings.TrimSpace(fmt.Sprintf("%+v", cfg.Extra)) == "{}" {
			cfg.Extra = c2.Extra
		} else if strings.TrimSpace(fmt.Sprintf("%+v", c2.Extra)) != "#cloud-config\n{}" {
			log.Warnf("More than 1 extra cloud-init section found, merging cloud-init is not supported and will be ignored from")
			cfg.Extra = c2.Extra
		}
	}
	if c2.Cleanup != nil {
		cfg.Cleanup = c2.Cleanup
	}
	cfg.Commands = append(cfg.Commands, c2.Commands...)
	cfg.PreCommands = append(cfg.PreCommands, c2.PreCommands...)
	cfg.PostCommands = append(cfg.PostCommands, c2.PostCommands...)
	cfg.Users = append(cfg.Users, c2.Users...)

	for k, v := range c2.Files {
		cfg.Files[k] = v
	}
	for k, v := range c2.Filesystem {
		cfg.Filesystem[k] = v
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
	pkgRepos := append(*cfg.PackageRepos, *c2.PackageRepos...)
	cfg.PackageRepos = &pkgRepos
	pkgs := append(*cfg.Packages, *c2.Packages...)
	cfg.Packages = &pkgs
	cfg.Timezone = c2.Timezone
	if c2.ContainerRuntime != nil {
		cfg.ContainerRuntime = c2.ContainerRuntime
	}
	cfg.TrustedCA = append(cfg.TrustedCA, c2.TrustedCA...)
	cfg.Limits = append(cfg.Limits, c2.Limits...)
	if c2.Kubernetes != nil {
		cfg.Kubernetes = c2.Kubernetes
	}
}
