package apps

import (
	"github.com/moshloop/konfigadm/pkg/phases"
	. "github.com/moshloop/konfigadm/pkg/types"
)

var Cleanup Phase = cleanup{}

type cleanup struct{}

func (c cleanup) ApplyPhase(sys *Config, ctx *SystemContext) ([]Command, Filesystem, error) {
	fs := Filesystem{}
	cmds := Commands{}

	if sys.Cleanup == nil || !*sys.Cleanup {
		return cmds.GetCommands(), fs, nil
	}

	for _, os := range phases.BaseOperatingSystems {
		cmds.AddAll(os.GetPackageManager().CleanupCaches().GetCommands()...)
	}

	cmds.Add("rm -rf /tmp/* || true").
		Add("rm -rf /var/log/cloud-init* || true").
		Add("rm -rf /usr/share/man/* || true").
		Add("rm -rf /usr/share/doc/* || true").
		Add("rm -f /etc/ssh/{ssh_host_dsa_key,ssh_host_dsa_key.pub,ssh_host_ecdsa_key,ssh_host_ecdsa_key.pub,ssh_host_ed25519_key,ssh_host_ed25519_key.pub,ssh_host_rsa_key,ssh_host_rsa_key.pub} || true").
		Add("find /var/cache -type f -exec rm -rf {} \\;").
		Add("find /var/log -type f | while read -r f; do echo -ne '' > \"$f\"; done;").
		Add("cloud-init clean || true")

	fs["/root/.bash_profile"] = File{Content: "shopt -s histappend"}
	fs["/etc/machine-id"] = File{Content: ""}
	fs["/root/.bash_history"] = File{Content: ""}
	fs["/root/.bashrc"] = File{Content: ""}

	// TODO
	// dd if=/dev/zero of=/EMPTY bs=1M  2>/dev/null || echo "dd exit code $? is suppressed"
	// rm -f /EMPTY
	//redhat
	//sed -i '/^\(HWADDR\|UUID\)=/d' /etc/sysconfig/network-scripts/ifcfg-*
	//debian /etc/udev/rules.d/70-persistent-net.rules

	return cmds.GetCommands(), fs, nil
}
