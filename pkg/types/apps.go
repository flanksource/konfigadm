package types

//ContainerRuntime installs a container runtime such as docker or CRI-O
type ContainerRuntime struct {
	Type    string `yaml:"type,omitempty"`
	Arg     string `yaml:"arg,omitempty"`
	Options string `yaml:"options,omitempty"`
	Version string `yaml:"version,omitempty"`
	//Images is a list of container images to pre-pull
	Images []string `yaml:"images,omitempty"`
}

//KubernetesSpec installs the packages and configures the system for kubernetes, it does not actually bootstrap and configure kubernetes itself
//Use kubeadm in a `command` to actually configure and start kubernetes
type KubernetesSpec struct {
	Version      string `yaml:"version,omitempty"`
	DownloadPath string `yaml:"download_path,omitempty"`
	ImagePrefix  string `yaml:"image_prefix,omitempty"`
}
