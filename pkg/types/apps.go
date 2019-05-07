package types

//ContainerRuntime installs a container runtime such as docker or CRI-O
type ContainerRuntime struct {
	Type    string `json:"type,omitempty"`
	Arg     string `json:"arg,omitempty"`
	Options string `json:"options,omitempty"`
	Version string `json:"version,omitempty"`
}

//KubernetesSpec installs the packages and configures the system for kubernetes, it does not actually bootstrap and configure kubernetes itself
//Use kubeadm in a `command` to actually configure and start kubernetes
type KubernetesSpec struct {
	Version      string `json:"version,omitempty"`
	DownloadPath string
	ImagePrefix  string
}
