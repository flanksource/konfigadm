package cmd

var images = map[string]string{
	"ubuntu:1804": "https://cloud-images.ubuntu.com/bionic/current/bionic-server-cloudimg-amd64.img",
	"ubuntu:1604": "https://cloud-images.ubuntu.com/xenial/current/xenial-server-cloudimg-amd64-disk1.img",
	"debian:9":    "https://cloud.debian.org/images/openstack/current-9/debian-9-openstack-amd64.qcow2",
	"centos:7":    "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.qcow2",
	"amazonLinux": "https://cdn.amazonlinux.com/os-images/2.0.20190612/kvm/amzn2-kvm-2.0.20190612-x86_64.xfs.gpt.qcow2",
	"fedora:30":   "https://download.fedoraproject.org/pub/fedora/linux/releases/30/Cloud/x86_64/images/Fedora-Cloud-Base-30-1.2.x86_64.qcow2",
	"fedora:29":   "https://download.fedoraproject.org/pub/fedora/linux/releases/29/Cloud/x86_64/images/Fedora-Cloud-Base-29-1.2.x86_64.qcow2",
}
