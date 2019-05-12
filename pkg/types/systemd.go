package types

import (
	. "github.com/moshloop/configadm/pkg/utils"
)

func (sys SystemD) ToUnitFile() string {
	return "[Unit]\n" + StructToIni(sys.Unit) + "\n" +
		"[Service]\n" + StructToIni(sys.Service) + "\n" +
		"[Install]\n" + StructToIni(sys.Install)
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
			Description: name,
		},
	}

}

// bool and int fields are modelled using interface{} to distinguish between nil (not provided) and empty values
// validation tags are used to enforce type

type SystemdUnit struct {
	After                       string
	AllowIsolate                interface{} `validate:"bool"`
	AssertACPower               string
	AssertArchitecture          string
	AssertCapability            string
	AssertDirectoryNotEmpty     string
	AssertFileIsExecutable      string
	AssertFileNotEmpty          string
	AssertFirstBoot             string
	AssertHost                  string
	AssertKernelCommandLine     string
	AssertNeedsUpdate           string
	AssertPathExists            string
	AssertPathExistsGlob        string
	AssertPathIsDirectory       string
	AssertPathIsMountPoint      string
	AssertPathIsReadWrite       string
	AssertPathIsSymbolicLink    string
	AssertSecurity              string
	AssertVirtualization        string
	Before                      string
	BindsTo                     string
	ConditionACPower            string
	ConditionArchitecture       string
	ConditionCapability         string
	ConditionDirectoryNotEmpty  string
	ConditionFileIsExecutable   string
	ConditionFileNotEmpty       string
	ConditionFirstBoot          string
	ConditionHost               string
	ConditionKernelCommandLine  string
	ConditionNeedsUpdate        string
	ConditionPathExists         string
	ConditionPathExistsGlob     string
	ConditionPathIsDirectory    string
	ConditionPathIsMountPoint   string
	ConditionPathIsReadWrite    string
	ConditionPathIsSymbolicLink string
	ConditionSecurity           string
	ConditionVirtualization     string
	Conflicts                   string
	DefaultDependencies         interface{} `validate:"bool"`
	Description                 string
	Documentation               string
	IgnoreOnIsolate             interface{} `validate:"bool"`
	JobTimeoutAction            string
	JobTimeoutRebootArgument    string
	JobTimeoutSec               string
	JoinsNamespaceOf            string
	//reboot-immediate, poweroff, poweroff-force or poweroff-immediates
	OnFailure             string
	OnFailureJobMode      string
	PartOf                string
	PropagatesReloadTo    string
	RebootArgument        string
	RefuseManualStart     interface{} `validate:"bool"`
	RefuseManualStop      interface{} `validate:"bool"`
	ReloadPropagatedFrom  string
	Requires              string
	RequiresMountsFor     string
	Requisite             string
	SourcePath            string
	StartLimitAction      string
	StartLimitBurst       string
	StartLimitIntervalSec string
	StopWhenUnneeded      interface{} `validate:"bool"`
	Wants                 string
}

type SystemdInstall struct {
	WantedBy        string
	RequiredBy      string
	DefaultInstance string
	Also            string
	Alias           string
}
type SystemdService struct {
	AmbientCapabilities      string
	AppArmorProfile          string
	CapabilityBoundingSet    string
	CPUAffinity              string
	CPUSchedulingPolicy      string
	CPUSchedulingPriority    string
	CPUSchedulingResetOnFork interface{} `validate:"bool"`
	DynamicUser              string
	Environment              string
	EnvironmentFile          string
	Group                    string
	IgnoreSIGPIPE            interface{} `validate:"bool"`
	InaccessiblePaths        string
	IOSchedulingClass        string
	IOSchedulingPriority     string
	LimitAS                  string
	LimitCORE                string
	LimitCPU                 string
	LimitDATA                string
	LimitFSIZE               string
	LimitLOCKS               string
	LimitMEMLOCK             string
	LimitMSGQUEUE            string
	LimitNICE                string
	LimitNOFILE              string
	LimitNPROC               string
	LimitRSS                 string
	LimitRTPRIO              string
	LimitRTTIME              string
	LimitSIGPENDING          string
	LimitSTACK               string
	MemoryDenyWriteExecute   interface{} `validate:"bool"`
	MountFlags               string
	Nice                     string
	NoNewPrivileges          interface{} `validate:"bool"`
	OOMScoreAdjust           string
	PAMName                  string
	PassEnvironment          string
	Personality              string
	PrivateDevices           interface{} `validate:"bool"`
	PrivateNetwork           interface{} `validate:"bool"`
	PrivateTmp               interface{} `validate:"bool"`
	PrivateUsers             interface{} `validate:"bool"`
	ProtectControlGroups     interface{} `validate:"bool"`
	ProtectHome              interface{} `validate:"bool"`
	ProtectKernelModules     string
	ProtectKernelTunables    interface{} `validate:"bool"`
	ProtectSystem            interface{} `validate:"bool"`
	ReadOnlyPaths            string
	ReadWritePaths           string
	RemoveIPC                interface{} `validate:"bool"`
	RestrictAddressFamilies  string
	RestrictNamespaces       interface{} `validate:"bool"`
	RestrictRealtime         interface{} `validate:"bool"`
	RootDirectory            string
	RuntimeDirectory         string
	RuntimeDirectoryMode     string
	SecureBits               string
	SELinuxContext           string
	SmackProcessLabel        string
	StandardError            string
	StandardInput            string
	StandardOutput           string
	SupplementaryGroups      string
	SyslogFacility           string
	SyslogIdentifier         string
	SyslogLevel              string
	SyslogLevelPrefix        interface{} `validate:"bool"`
	SystemCallArchitectures  string
	SystemCallErrorNumber    string
	SystemCallFilter         string
	TimerSlackNSec           string
	TTYPath                  string
	TTYReset                 string
	TTYVHangup               string
	TTYVTDisallocate         string
	Umask                    string
	User                     string
	UtmpIdentifier           string
	UtmpMode                 string
	WorkingDirectory         string //Takes a D-Bus bus name that this service is reachable as. This option is mandatory for services where Type= is set to dbus.
	BusName                  string
	ExecReload               string
	ExecStart                string
	ExecStartPost            string
	ExecStartPre             string
	ExecStop                 string
	ExecStopPost             string
	FailureAction            string
	FileDescriptorStoreMax   interface{} `validate:"int"`
	GuessMainPID             interface{} `validate:"bool"`
	NonBlocking              interface{} `validate:"bool"`
	NotifyAccess             string
	PermissionsStartOnly     interface{} `validate:"int"`
	PIDFile                  string
	RemainAfterExit          interface{} `validate:"bool"`
	//	always,	on-success,	on-failure,	on-abnormal,	on-abort,	on-watchdog
	Restart                  string
	RestartForceExitStatus   interface{} `validate:"int"`
	RestartPreventExitStatus interface{} `validate:"int"`
	//Configures the time to sleep before restarting a service (as configured with Restart=). Takes a unit-less value in seconds, or a time span value such as "5min 20s". Defaults to 100ms.
	RestartSec             string
	RootDirectoryStartOnly interface{} `validate:"int"`
	RuntimeMaxSec          string
	SuccessExitStatus      interface{} `validate:"int"`
	TimeoutSec             string
	TimeoutStartSec        string
	TimeoutStopSec         string
	//simple, forking, oneshot, dbus, notify or idle
	Type        string
	WatchdogSec string
}

type SystemD struct {
	Install SystemdInstall
	Service SystemdService
	Unit    SystemdUnit
}
