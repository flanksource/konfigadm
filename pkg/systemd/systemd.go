package systemd

import (
	. "github.com/moshloop/cloud-config/pkg/utils"
)

func (sys SystemD) ToUnitFile() string {
	return "[Unit]\n" + StructToIni(sys.Unit) + "\n[Service]\n" + StructToIni(sys.Service) + "\n[Install]\n" + StructToIni(sys.Install)
}

func DefaultSystemdService(name string) SystemD {
	return SystemD{
		Install: SystemdInstall{
			WantedBy: "multi-user.target",
		},
		Service: SystemdService{
			Restart:    "on-failure",
			RestartSec: "60",
		},
		Unit: SystemdUnit{
			StopWhenUnneeded: true,
			Description:      name,
		},
	}

}
