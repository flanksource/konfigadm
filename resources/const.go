package resources

var ContainerdService = FSMustString(false, "/containerd.service")
var KubeletConf = FSMustString(false, "/kubeadm.service")
var VfsStorageDriver = FSMustString(false, "/daemon.json")
