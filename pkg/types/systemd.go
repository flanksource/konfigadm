package types

import (
	. "github.com/flanksource/konfigadm/pkg/utils" // nolint: golint, stylecheck
)

//Service is a systemd service to be installed and started
type Service struct {
	Name        string            `yaml:"name,omitempty"`
	ExecStart   string            `yaml:"exec_start,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Extra       SystemD           `yaml:"extra,omitempty"`
	// TODO: capabilities
}

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
			RestartSec: "10",
		},
		Unit: SystemdUnit{
			Description: name,
		},
	}
}

// bool and int fields are modelled using interface{} to distinguish between nil (not provided) and empty values
// validation tags are used to enforce type

type SystemdUnit struct {
	After                       string      `yaml:"after,omitempty"`
	AllowIsolate                interface{} `validate:"bool" yaml:"allow_isolate,omitempty"`
	AssertACPower               string      `yaml:"assert_ac_power,omitempty"`
	AssertArchitecture          string      `yaml:"assert_architecture,omitempty"`
	AssertCapability            string      `yaml:"assert_capability,omitempty"`
	AssertDirectoryNotEmpty     string      `yaml:"assert_directory_not_empty,omitempty"`
	AssertFileIsExecutable      string      `yaml:"assert_file_is_executable,omitempty"`
	AssertFileNotEmpty          string      `yaml:"assert_file_not_empty,omitempty"`
	AssertFirstBoot             string      `yaml:"assert_first_boot,omitempty"`
	AssertHost                  string      `yaml:"assert_host,omitempty"`
	AssertKernelCommandLine     string      `yaml:"assert_kernel_command_line,omitempty"`
	AssertNeedsUpdate           string      `yaml:"assert_needs_update,omitempty"`
	AssertPathExists            string      `yaml:"assert_path_exists,omitempty"`
	AssertPathExistsGlob        string      `yaml:"assert_path_exists_glob,omitempty"`
	AssertPathIsDirectory       string      `yaml:"assert_path_is_directory,omitempty"`
	AssertPathIsMountPoint      string      `yaml:"assert_path_is_mount_point,omitempty"`
	AssertPathIsReadWrite       string      `yaml:"assert_path_is_read_write,omitempty"`
	AssertPathIsSymbolicLink    string      `yaml:"assert_path_is_symbolic_link,omitempty"`
	AssertSecurity              string      `yaml:"assert_security,omitempty"`
	AssertVirtualization        string      `yaml:"assert_virtualization,omitempty"`
	Before                      string      `yaml:"before,omitempty"`
	BindsTo                     string      `yaml:"binds_to,omitempty"`
	ConditionACPower            string      `yaml:"condition_ac_power,omitempty"`
	ConditionArchitecture       string      `yaml:"condition_architecture,omitempty"`
	ConditionCapability         string      `yaml:"condition_capability,omitempty"`
	ConditionDirectoryNotEmpty  string      `yaml:"condition_directory_not_empty,omitempty"`
	ConditionFileIsExecutable   string      `yaml:"condition_file_is_executable,omitempty"`
	ConditionFileNotEmpty       string      `yaml:"condition_file_not_empty,omitempty"`
	ConditionFirstBoot          string      `yaml:"condition_first_boot,omitempty"`
	ConditionHost               string      `yaml:"condition_host,omitempty"`
	ConditionKernelCommandLine  string      `yaml:"condition_kernel_command_line,omitempty"`
	ConditionNeedsUpdate        string      `yaml:"condition_needs_update,omitempty"`
	ConditionPathExists         string      `yaml:"condition_path_exists,omitempty"`
	ConditionPathExistsGlob     string      `yaml:"condition_path_exists_glob,omitempty"`
	ConditionPathIsDirectory    string      `yaml:"condition_path_is_directory,omitempty"`
	ConditionPathIsMountPoint   string      `yaml:"condition_path_is_mount_point,omitempty"`
	ConditionPathIsReadWrite    string      `yaml:"condition_path_is_read_write,omitempty"`
	ConditionPathIsSymbolicLink string      `yaml:"condition_path_is_symbolic_link,omitempty"`
	ConditionSecurity           string      `yaml:"condition_security,omitempty"`
	ConditionVirtualization     string      `yaml:"condition_virtualization,omitempty"`
	Conflicts                   string      `yaml:"conflicts,omitempty"`
	DefaultDependencies         interface{} `validate:"bool" yaml:"default_dependencies,omitempty"`
	Description                 string      `yaml:"description,omitempty"`
	Documentation               string      `yaml:"documentation,omitempty"`
	IgnoreOnIsolate             interface{} `validate:"bool" yaml:"ignore_on_isolate,omitempty"`
	JobTimeoutAction            string      `yaml:"job_timeout_action,omitempty"`
	JobTimeoutRebootArgument    string      `yaml:"job_timeout_reboot_argument,omitempty"`
	JobTimeoutSec               string      `yaml:"job_timeout_sec,omitempty"`
	JoinsNamespaceOf            string      `yaml:"joins_namespace_of,omitempty"`
	//reboot-immediate, poweroff, poweroff-force or poweroff-immediates
	OnFailure             string      `yaml:"on_failure,omitempty"`
	OnFailureJobMode      string      `yaml:"on_failure_job_mode,omitempty"`
	PartOf                string      `yaml:"part_of,omitempty"`
	PropagatesReloadTo    string      `yaml:"propagates_reload_to,omitempty"`
	RebootArgument        string      `yaml:"reboot_argument,omitempty"`
	RefuseManualStart     interface{} `validate:"bool" yaml:"refuse_manual_start,omitempty"`
	RefuseManualStop      interface{} `validate:"bool" yaml:"refuse_manual_stop,omitempty"`
	ReloadPropagatedFrom  string      `yaml:"reload_propagated_from,omitempty"`
	Requires              string      `yaml:"requires,omitempty"`
	RequiresMountsFor     string      `yaml:"requires_mounts_for,omitempty"`
	Requisite             string      `yaml:"requisite,omitempty"`
	SourcePath            string      `yaml:"source_path,omitempty"`
	StartLimitAction      string      `yaml:"start_limit_action,omitempty"`
	StartLimitBurst       string      `yaml:"start_limit_burst,omitempty"`
	StartLimitIntervalSec string      `yaml:"start_limit_interval_sec,omitempty"`
	StopWhenUnneeded      interface{} `validate:"bool" yaml:"stop_when_unneeded,omitempty"`
	Wants                 string      `yaml:"wants,omitempty"`
}

type SystemdInstall struct {
	WantedBy        string `yaml:"wanted_by,omitempty"`
	RequiredBy      string `yaml:"required_by,omitempty"`
	DefaultInstance string `yaml:"default_instance,omitempty"`
	Also            string `yaml:"also,omitempty"`
	Alias           string `yaml:"alias,omitempty"`
}
type SystemdService struct {
	AmbientCapabilities      string      `yaml:"ambient_capabilities,omitempty"`
	AppArmorProfile          string      `yaml:"app_armor_profile,omitempty"`
	CapabilityBoundingSet    string      `yaml:"capability_bounding_set,omitempty"`
	CPUAffinity              string      `yaml:"cpu_affinity,omitempty"`
	CPUSchedulingPolicy      string      `yaml:"cpu_scheduling_policy,omitempty"`
	CPUSchedulingPriority    string      `yaml:"cpu_scheduling_priority,omitempty"`
	CPUSchedulingResetOnFork interface{} `validate:"bool" yaml:"cpu_scheduling_reset_on_fork,omitempty"`
	DynamicUser              string      `yaml:"dynamic_user,omitempty"`
	Environment              string      `yaml:"environment,omitempty"`
	EnvironmentFile          string      `yaml:"environment_file,omitempty"`
	Group                    string      `yaml:"group,omitempty"`
	IgnoreSIGPIPE            interface{} `validate:"bool" yaml:"ignore_sigpipe,omitempty"`
	InaccessiblePaths        string      `yaml:"inaccessible_paths,omitempty"`
	IOSchedulingClass        string      `yaml:"io_scheduling_class,omitempty"`
	IOSchedulingPriority     string      `yaml:"io_scheduling_priority,omitempty"`
	LimitAS                  string      `yaml:"limit_as,omitempty"`
	LimitCORE                string      `yaml:"limit_core,omitempty"`
	LimitCPU                 string      `yaml:"limit_cpu,omitempty"`
	LimitDATA                string      `yaml:"limit_data,omitempty"`
	LimitFSIZE               string      `yaml:"limit_fsize,omitempty"`
	LimitLOCKS               string      `yaml:"limit_locks,omitempty"`
	LimitMEMLOCK             string      `yaml:"limit_memlock,omitempty"`
	LimitMSGQUEUE            string      `yaml:"limit_msgqueue,omitempty"`
	LimitNICE                string      `yaml:"limit_nice,omitempty"`
	LimitNOFILE              string      `yaml:"limit_nofile,omitempty"`
	LimitNPROC               string      `yaml:"limit_nproc,omitempty"`
	LimitRSS                 string      `yaml:"limit_rss,omitempty"`
	LimitRTPRIO              string      `yaml:"limit_rtprio,omitempty"`
	LimitRTTIME              string      `yaml:"limit_rttime,omitempty"`
	LimitSIGPENDING          string      `yaml:"limit_sigpending,omitempty"`
	LimitSTACK               string      `yaml:"limit_stack,omitempty"`
	MemoryDenyWriteExecute   interface{} `validate:"bool" yaml:"memory_deny_write_execute,omitempty"`
	MountFlags               string      `yaml:"mount_flags,omitempty"`
	Nice                     string      `yaml:"nice,omitempty"`
	NoNewPrivileges          interface{} `validate:"bool" yaml:"no_new_privileges,omitempty"`
	OOMScoreAdjust           string      `yaml:"oom_score_adjust,omitempty"`
	PAMName                  string      `yaml:"pam_name,omitempty"`
	PassEnvironment          string      `yaml:"pass_environment,omitempty"`
	Personality              string      `yaml:"personality,omitempty"`
	PrivateDevices           interface{} `validate:"bool" yaml:"private_devices,omitempty"`
	PrivateNetwork           interface{} `validate:"bool" yaml:"private_network,omitempty"`
	PrivateTmp               interface{} `validate:"bool" yaml:"private_tmp,omitempty"`
	PrivateUsers             interface{} `validate:"bool" yaml:"private_users,omitempty"`
	ProtectControlGroups     interface{} `validate:"bool" yaml:"protect_control_groups,omitempty"`
	ProtectHome              interface{} `validate:"bool" yaml:"protect_home,omitempty"`
	ProtectKernelModules     string      `yaml:"protect_kernel_modules,omitempty"`
	ProtectKernelTunables    interface{} `validate:"bool" yaml:"protect_kernel_tunables,omitempty"`
	ProtectSystem            interface{} `validate:"bool" yaml:"protect_system,omitempty"`
	ReadOnlyPaths            string      `yaml:"read_only_paths,omitempty"`
	ReadWritePaths           string      `yaml:"read_write_paths,omitempty"`
	RemoveIPC                interface{} `validate:"bool" yaml:"remove_ipc,omitempty"`
	RestrictAddressFamilies  string      `yaml:"restrict_address_families,omitempty"`
	RestrictNamespaces       interface{} `validate:"bool" yaml:"restrict_namespaces,omitempty"`
	RestrictRealtime         interface{} `validate:"bool" yaml:"restrict_realtime,omitempty"`
	RootDirectory            string      `yaml:"root_directory,omitempty"`
	RuntimeDirectory         string      `yaml:"runtime_directory,omitempty"`
	RuntimeDirectoryMode     string      `yaml:"runtime_directory_mode,omitempty"`
	SecureBits               string      `yaml:"secure_bits,omitempty"`
	SELinuxContext           string      `yaml:"se_linux_context,omitempty"`
	SmackProcessLabel        string      `yaml:"smack_process_label,omitempty"`
	StandardError            string      `yaml:"standard_error,omitempty"`
	StandardInput            string      `yaml:"standard_input,omitempty"`
	StandardOutput           string      `yaml:"standard_output,omitempty"`
	SupplementaryGroups      string      `yaml:"supplementary_groups,omitempty"`
	SyslogFacility           string      `yaml:"syslog_facility,omitempty"`
	SyslogIdentifier         string      `yaml:"syslog_identifier,omitempty"`
	SyslogLevel              string      `yaml:"syslog_level,omitempty"`
	SyslogLevelPrefix        interface{} `validate:"bool" yaml:"syslog_level_prefix,omitempty"`
	SystemCallArchitectures  string      `yaml:"system_call_architectures,omitempty"`
	SystemCallErrorNumber    string      `yaml:"system_call_error_number,omitempty"`
	SystemCallFilter         string      `yaml:"system_call_filter,omitempty"`
	TimerSlackNSec           string      `yaml:"timer_slack_n_sec,omitempty"`
	TTYPath                  string      `yaml:"tty_path,omitempty"`
	TTYReset                 string      `yaml:"tty_reset,omitempty"`
	TTYVHangup               string      `yaml:"ttyv_hangup,omitempty"`
	TTYVTDisallocate         string      `yaml:"ttyvt_disallocate,omitempty"`
	Umask                    string      `yaml:"umask,omitempty"`
	User                     string      `yaml:"user,omitempty"`
	UtmpIdentifier           string      `yaml:"utmp_identifier,omitempty"`
	UtmpMode                 string      `yaml:"utmp_mode,omitempty"`
	WorkingDirectory         string      `yaml:"working_directory,omitempty"` //Takes a D-Bus bus name that this service is reachable as. This option is mandatory for services where Type= is set to dbus.
	BusName                  string      `yaml:"bus_name,omitempty"`
	ExecReload               string      `yaml:"exec_reload,omitempty"`
	ExecStart                string      `yaml:"exec_start,omitempty"`
	ExecStartPost            string      `yaml:"exec_start_post,omitempty"`
	ExecStartPre             string      `yaml:"exec_start_pre,omitempty"`
	ExecStop                 string      `yaml:"exec_stop,omitempty"`
	ExecStopPost             string      `yaml:"exec_stop_post,omitempty"`
	FailureAction            string      `yaml:"failure_action,omitempty"`
	FileDescriptorStoreMax   interface{} `validate:"int" yaml:"file_descriptor_store_max,omitempty"`
	GuessMainPID             interface{} `validate:"bool" yaml:"guess_main_pid,omitempty"`
	NonBlocking              interface{} `validate:"bool" yaml:"non_blocking,omitempty"`
	NotifyAccess             string      `yaml:"notify_access,omitempty"`
	PermissionsStartOnly     interface{} `validate:"int" yaml:"permissions_start_only,omitempty"`
	PIDFile                  string      `yaml:"pid_file,omitempty"`
	RemainAfterExit          interface{} `validate:"bool" yaml:"remain_after_exit,omitempty"`
	//	always,	on-success,	on-failure,	on-abnormal,	on-abort,	on-watchdog
	Restart                  string      `yaml:"restart,omitempty"`
	RestartForceExitStatus   interface{} `validate:"int" yaml:"restart_force_exit_status,omitempty"`
	RestartPreventExitStatus interface{} `validate:"int" yaml:"restart_prevent_exit_status,omitempty"`
	//Configures the time to sleep before restarting a service (as configured with Restart=). Takes a unit-less value in seconds, or a time span value such as "5min 20s". Defaults to 100ms.
	RestartSec             string      `yaml:"restart_sec,omitempty"`
	RootDirectoryStartOnly interface{} `validate:"int" yaml:"root_directory_start_only,omitempty"`
	RuntimeMaxSec          string      `yaml:"runtime_max_sec,omitempty"`
	SuccessExitStatus      interface{} `validate:"int" yaml:"success_exit_status,omitempty"`
	TimeoutSec             string      `yaml:"timeout_sec,omitempty"`
	TimeoutStartSec        string      `yaml:"timeout_start_sec,omitempty"`
	TimeoutStopSec         string      `yaml:"timeout_stop_sec,omitempty"`
	//simple, forking, oneshot, dbus, notify or idle
	Type        string `yaml:"type,omitempty"`
	WatchdogSec string `yaml:"watchdog_sec,omitempty"`
}

type SystemD struct {
	Install SystemdInstall
	Service SystemdService
	Unit    SystemdUnit
}
