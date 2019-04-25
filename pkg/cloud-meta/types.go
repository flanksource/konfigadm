package cloudmeta

import (
	cloudinit "github.com/moshloop/configadm/pkg/cloud-init"
	. "github.com/moshloop/configadm/pkg/systemd"
)

type Port struct {
	Port   int `json:"port,omitempty"  validate:"min=1,max=65536"`
	Target int `json:"target,omitempty"  validate:"min=1,max=65536"`
}

type Container struct {

	//The name of the service (e.g systemd unit name or deployment name)
	Service string `json:"service,omitempty"`
	//A map of environment variables to pass through
	Env map[string]string `json:"env,omitempty"`
	//A map of labels to add to the container
	Labels map[string]string `json:"labels,omitempty"`
	//Additional arguments to the docker run command e.g. -p 8080:8080
	DockerOpts string `json:"docker_opts,omitempty"`
	//Additional options to the docker client e.g. -H unix:///tmp/var/run/docker.sock
	DockerClientArgs string `json:"docker_client_args,omitempty"`
	//Additional arguments to the container
	Args     string   `json:"args,omitempty"`
	Ports    []Port   `json:"ports,omitempty"`
	Commands []string `json:"commands,omitempty"`
	//Map of files to mount into the container
	Files map[string]string `json:"files,omitempty"`
	//Map of templates to mount into the container
	Templates map[string]string `json:"templates,omitempty"`
	Volumes   map[string]string `json:"volumes,omitempty"`
	//CPU limit in cores (Defaults to 1 )
	Cpu int `json:"cpu,omitempty" validate:"min=0,max=32"`
	//	Memory Limit in MB. (Defaults to 1024)
	Mem int `json:"mem,omitempty" validate:"min=0,max=1048576"`
	//default:	user-bridge	 only
	Network string `json:"network,omitempty"`
	// default: 1
	Replicas int `json:"replicas,omitempty"`
}

type ContainerRuntime struct {
	Type    string `json:"type,omitempty"`
	Arg     string `json:"arg,omitempty"`
	Options string `json:"options,omitempty"`
	Version string `json:"version,omitempty"`
}

type Kubernetes struct {
	Version string `json:"version,omitempty"`
}

type Service struct {
	ExecStart   string            `json:"exec_start,omitempty"`
	Environment map[string]string `json:"environment,omitempty"`
	Extra       SystemD           `json:"extra,omitempty"`
}

type SystemConfig struct {
	Commands         []string           `json:"commands,omitempty"`
	PreCommands      []string           `json:"pre_commands,omitempty"`
	PostCommands     []string           `json:"post_commands,omitempty"`
	Files            map[string]string  `json:"files,omitempty"`
	Templates        map[string]string  `json:"templates,omitempty"`
	Sysctls          map[string]string  `json:"sysctls,omitempty"`
	Packages         []string           `json:"packages,omitempty"`
	PackageRepos     []string           `json:"package_repos,omitempty"`
	Images           []string           `json:"images,omitempty"`
	Containers       []Container        `json:"containers,omitempty"`
	ContainerRuntime ContainerRuntime   `json:"container_runtime,omitempty"`
	Kubernetes       Kubernetes         `json:"kubernetes,omitempty"`
	Environment      map[string]string  `json:"environment,omitempty"`
	Timezone         string             `json:"timezone,omitempty"`
	Extra            CloudInit          `json:"extra,omitempty"`
	Services         map[string]Service `json:"services,omitempty"`
	Users            []User             `json:"users,omitempty"`
	Context          SystemContext
}

type SystemContext struct {
	Vars map[string]string
	Name string
}
