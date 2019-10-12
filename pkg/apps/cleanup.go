package apps

import (
	"github.com/moshloop/konfigadm/pkg/build"
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
		cmds.AddAll(os.GetPackageManager().CleanupCaches().WithTags(os.GetTags()...).GetCommands()...)
	}

	if ctx.CaptureLogs != "" {
		for _, cmd := range build.CaptureLogCommands() {
			cmds.Add(cmd)
		}
	}
	cmds.
		Add("rm -rf /tmp/* || true").
		Add("rm -rf /var/tmp/* || true").
		Add("rm -rf /usr/share/man/* || true").
		Add("rm -rf /usr/share/doc/* || true").
		Add("rm -rf /var/lib/dhclient/* || true").
		Add("rm /etc/netplan/50-cloud-init.yaml || true").
		Add("rm /etc/udev/rules.d/70-persistent-net.rules || true").
		Add("rm -f /etc/ssh/ssh_host_* || true").
		Add("sed -i '/^\\(HWADDR\\|UUID\\)=/d' /etc/sysconfig/network-scripts/ifcfg-* || true").
		Add("find /var/cache -type f -exec rm -rf {} \\;").
		Add("find /var/log -type f | while read -r f; do echo -ne '' > \"$f\"; done;").
		Add("rm -rf /var/run/cloud-init || true").
		Add("rm -rf /var/lib/cloud || true").
		Add("journalctl --rotate && sleep 5 && journalctl --vacuum-time=1s").
		Add("export MACHINE_ID=$(cat /etc/machine-id)").
		Add("echo -ne > /etc/machine-id").
		Add("[[ -e /var/lib/dbus/machine-id ]] && echo -ne > /var/lib/dbus/machine-id").
		Add("echo -ne > /root/.bash_history").
		Add("echo Finished cleanup on $(date) with machine-id: $(cat /etc/machine-id) old: $MACHINE_ID > /var/log/cleanup.log").
		Add("dd if=/dev/zero of=/EMPTY bs=1M  2>/dev/null || true;  rm -f /EMPTY")

	return cmds.GetCommands(), fs, nil
}
