package cmd

import (
	"github.com/flanksource/konfigadm/pkg/types"
	"github.com/flanksource/konfigadm/pkg/utils"
)

var images = map[string]Image{
	"ubuntu2004": Image{
		Url:            "https://cloud-images.ubuntu.com/releases/20.04/release-{{.version}}/ubuntu-20.04-server-cloudimg-amd64.img",
		DefaultVersion: "20200720",
		Tags:           []types.Flag{types.UBUNTU, types.DEBIAN_LIKE},
	},
	"ubuntu1804": Image{
		Url:            "https://cloud-images.ubuntu.com/releases/18.04/release-{{.version}}/ubuntu-18.04-server-cloudimg-amd64.img",
		DefaultVersion: "20190617",
		Tags:           []types.Flag{types.UBUNTU, types.DEBIAN_LIKE},
	},
	"ubuntu1604": Image{
		Url:            "https://cloud-images.ubuntu.com/releases/16.04/release-{{.version}}/ubuntu-16.04-server-cloudimg-amd64-disk1.img",
		DefaultVersion: "20190628",
		Tags:           []types.Flag{types.UBUNTU, types.DEBIAN_LIKE},
	},
	"debian": Image{
		Url:            "https://cloud.debian.org/images/openstack/archive/{{.version}}/debian-{{.version}}-openstack-amd64.qcow2",
		DefaultVersion: "9.9.3-20190618",
		Tags:           []types.Flag{types.DEBIAN, types.DEBIAN_LIKE},
	},
	"centos7": Image{
		Url:            "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud-{{.version}}.qcow2",
		DefaultVersion: "1905",
		Tags:           []types.Flag{types.CENTOS, types.REDHAT_LIKE},
	},
	"amazonLinux": Image{
		Url:            "https://cdn.amazonlinux.com/os-images/2.0.20190612/kvm/amzn2-kvm-2.0.{{.version}}-x86_64.xfs.gpt.qcow2",
		DefaultVersion: "20190612",
		Tags:           []types.Flag{types.AMAZON_LINUX, types.REDHAT_LIKE},
	},
	"fedora": Image{
		Url:            "https://download.fedoraproject.org/pub/fedora/linux/releases/{{.version}}/Cloud/x86_64/images/Fedora-Cloud-Base-{{.version}}-1.2.x86_64.qcow2",
		DefaultVersion: "30",
		Tags:           []types.Flag{types.FEDORA},
	},
}

type Image struct {
	Alias          string
	Url            string
	Tags           []types.Flag
	Version        string
	DefaultVersion string
}

func (i Image) GetURL() string {
	vars := map[string]string{"version": i.DefaultVersion}
	if i.Version != "" {
		vars["version"] = i.Version
	}
	return utils.Interpolate(i.Url, vars)
}
