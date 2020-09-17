package types

import (
	"fmt"
	"strings"

	"gopkg.in/flanksource/yaml.v3"
)

//Package includes the package name, modifiers (mark, uninstall) and runtime tags
type Package struct {
	Name      string
	Version   string
	Mark      bool
	Uninstall bool
	Flags     []Flag
}

func CompareVersions(version string, compareTo string) bool {
	if strings.Contains(compareTo, "==") {
		compareTo = strings.Split(compareTo, "==")[1]
	} else if strings.Contains(compareTo, "=") {
		compareTo = strings.Split(compareTo, "=")[1]
	}
	return version == compareTo
}

//PackageRepo includes the URL for a package repo, GPG key (if applicable) and runtime tags
type PackageRepo struct {
	Name            string            `yaml:"name,omitempty"`
	URL             string            `yaml:"url,omitempty"`
	GPGKey          string            `yaml:"gpgKey,omitempty"`
	Channel         string            `yaml:"channel,omitempty"`
	VersionCodeName string            `yaml:"versionCodeName,omitempty"`
	Flags           []Flag            `yaml:"tags,omitempty"`
	ExtraArgs       map[string]string `yaml:"extraArgs,omitempty"`
}

func (p PackageRepo) String() string {
	return p.Name
}

type PackageManager interface {
	Install(pkg ...string) Commands
	Uninstall(pkg ...string) Commands
	Mark(pkg ...string) Commands
	AddRepo(url string, channel string, versionCodeName string, name string, gpgKey string, extraArgs map[string]string) Commands
	GetInstalledVersion(pkg string) string
	CleanupCaches() Commands
	Update() Commands
	GetKernelPackageNames(version string) (string, string)
}

func (p Package) String() string {
	return p.Name + fmt.Sprintf("%s", p.Flags)
}

func (p Package) VersionedName() string {
	if p.Version == "" {
		return p.Name
	}
	return p.Name + "=" + p.Version
}

//AppendPackages is a helper function to add new packages
func (cfg *Config) AppendPackages(flag *Flag, packages ...Package) *Config {
	for _, pkg := range packages {
		if flag != nil {
			pkg.Flags = []Flag{*flag}
		}
		pkgs := append(*cfg.Packages, pkg)
		cfg.Packages = &pkgs
	}
	return cfg
}

//AddPackage is a helper function to add new packages
func (cfg *Config) AddPackage(names string, flag *Flag) *Config {
	for _, name := range strings.Split(names, " ") {
		pkg := Package{}
		parsePackage(name, &pkg)
		if flag != nil {
			pkg.Flags = []Flag{*flag}
		}
		pkgs := append(*cfg.Packages, pkg)
		cfg.Packages = &pkgs
	}
	return cfg
}

//AddTarPackage is a helper function to add new packages
func (cfg *Config) AddTarPackage(pkg TarPackage, flags ...Flag) *Config {
	pkg.Flags = flags
	cfg.TarPackages = append(cfg.TarPackages, pkg)
	return cfg
}

//AddPackageRepo is a helper function to add new packages repos
func (cfg *Config) AddPackageRepo(url string, gpg string, flag Flag) *Config {
	pkg := PackageRepo{
		URL: url,
	}

	if gpg != "" {
		pkg.GPGKey = gpg
	}
	return cfg.AppendPackageRepo(pkg, flag)
}

//AppendPackageRepo appends a new package repository to the list
func (cfg *Config) AppendPackageRepo(repo PackageRepo, flags ...Flag) *Config {
	repo.Flags = append(repo.Flags, flags...)
	if repo.Channel == "" {
		repo.Channel = "main"
	}
	pkgs := append(*cfg.PackageRepos, repo)
	cfg.PackageRepos = &pkgs
	return cfg
}

//MarshalYAML adds tags as comments
func (p Package) MarshalYAML() (interface{}, error) {
	return &yaml.Node{
		Kind:        yaml.ScalarNode,
		Tag:         "!!str",
		LineComment: Marshall(p.Flags),
		Value:       p.Name,
	}, nil
}

func parsePackage(name string, p *Package) *Package {
	if strings.HasPrefix(name, "!") {
		p.Name = name[1:]
		p.Uninstall = true
	} else if strings.HasPrefix(name, "+") {
		p.Name = name[1:]
		p.Mark = true
	} else {
		p.Name = name
	}

	if strings.Contains(name, "=") {
		parts := strings.Split(name, "=")
		p.Name = parts[0]
		p.Version = parts[len(parts)-1]
	}
	return p
}

//UnmarshalYAML decodes comments into tags and parses modifiers for packages
func (p *Package) UnmarshalYAML(node *yaml.Node) error {
	parsePackage(node.Value, p)

	comment := node.LineComment
	if !strings.Contains(comment, "#") {
		return nil
	}
	comment = comment[1:]
	for _, flag := range strings.Split(comment, " ") {
		if FLAG, ok := FLAG_MAP[flag]; ok {
			p.Flags = append(p.Flags, FLAG)
		} else {
			return fmt.Errorf("Unknown flag: %s", flag)
		}
	}
	return nil
}
