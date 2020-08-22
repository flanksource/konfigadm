package cloudinit

type AptSource struct {
	// a sources.list entry
	// sources can use $MIRROR, $PRIMARY, $SECURITY and $RELEASE replacement variables.
	// They will be replaced with the default or specified mirrors and the running release.
	// The entry below would be possibly turned into:
	//   source: deb http://archive.ubuntu.com/ubuntu xenial multiverse
	Source string `yaml:"source,omitempty"`
	//Importing a gpg key for a given key id. Used keyserver defaults to keyserver.ubuntu.com
	Keyid string `yaml:"keyid,omitempty"`
	// specify an alternate keyserver to pull keys from that were specified by keyid
	KeyServer string `yaml:"keyserver,omitempty"`
	//providing a raw PGP key
	Key      string `yaml:"key,omitempty"`
	Filename string `yaml:"filename,omitempty"`
}
type Apt struct {
	// Preserves the existing /etc/apt/sources.list
	// Default: false - do overwrite sources_list. If set to true then any
	// "mirrors" configuration will have no effect.
	// Set to true to avoid affecting sources.list. In that case only
	// "extra" source specifications will be written into
	// /etc/apt/sources.list.d/*
	PreserveSourcesList bool `yaml:"preserve_sources_list,omitempty"`
	// If given, those suites are removed from sources.list after all other
	// modifications have been made.
	// Suites are even disabled if no other modification was made,
	// but not if is preserve_sources_list is active.
	// There is a special alias "$RELEASE" as in the sources that will be replace
	// by the matching release.
	//
	// To ease configuration and improve readability the following common ubuntu
	// suites will be automatically mapped to their full definition.
	// updates   => $RELEASE-updates
	// backports => $RELEASE-backports
	// security  => $RELEASE-security
	// proposed  => $RELEASE-proposed
	// release   => $RELEASE
	//
	// There is no harm in specifying a suite to be disabled that is not found in
	// the source.list file (just a no-op then)
	//
	// Note: Lines don't get deleted, but disabled by being converted to a comment.
	// The following example disables all usual defaults except $RELEASE-security.
	// On top it disables a custom suite called "mysuite"
	DisableSuites []string `yaml:"disable_suites,omitempty"`
	// Default: none - instead it is auto select based on cloud metadata
	// so if neither "uri" nor "search", nor "search_dns" is set (the default)
	// then use the mirror provided by the DataSource found.
	// In EC2, that means using <region>.ec2.archive.ubuntu.com
	//
	// define a custom (e.g. localized) mirror that will be used in sources.list
	// and any custom sources entries for deb / deb-src lines.
	//
	// One can set primary and security mirror to different uri's
	// the child elements to the keys primary and secondary are equivalent
	//
	// If multiple of a category are given
	//   1. uri
	//   2. search
	//   3. search_dns
	// the first defining a valid mirror wins (in the order as defined here,
	// not the order as listed in the config).
	Primary []struct {
		// arches is list of architectures the following config applies to
		// the special keyword "default" applies to any architecture not explicitly
		// listed.
		Arches []string `yaml:"arches,omitempty"`
		URI    string   `yaml:"uri,omitempty"`
		// via search one can define lists that are tried one by one.
		// The first with a working DNS resolution (or if it is an IP) will be
		// picked. That way one can keep one configuration for multiple
		// subenvironments that select the working one.
		Search []string `yaml:"search,omitempty"`
		// if no mirror is provided by uri or search but 'search_dns' is
		// true, then search for dns names '<distro>-mirror' in each of
		// - fqdn of this host per cloud metadata
		// - localdomain
		// - no domain (which would search domains listed in /etc/resolv.conf)
		// If there is a dns entry for <distro>-mirror, then it is assumed that
		// there is a distro mirror at http://<distro>-mirror.<domain>/<distro>
		//
		// That gives the cloud provider the opportunity to set mirrors of a distro
		// up and expose them only by creating dns entries.
		//
		// if none of that is found, then the default distro mirror is used
		SearchDNS bool `yaml:"search_dns,omitempty"`
	} `yaml:"primary,omitempty"`
	Security struct {
		URI string `yaml:"uri,omitempty"`
	} `yaml:"security,omitempty"`
	// Provide a custom template for rendering sources.list
	// without one provided cloud-init uses builtin templates for
	// ubuntu and debian.
	// Within these sources.list templates you can use the following replacement
	// variables (all have sane Ubuntu defaults, but mirrors can be overwritten
	// as needed (see above)):
	// => $RELEASE, $MIRROR, $PRIMARY, $SECURITY
	SourcesList string `yaml:"sources_list,omitempty"`

	// Any apt config string that will be made available to apt
	// see the APT.CONF(5) man page for details what can be specified
	Conf string `yaml:"conf,omitempty"`
	// Proxies are the most common apt.conf option, so that for simplified use
	// there is a shortcut for those. Those get automatically translated into the
	// correct Acquire::*::Proxy statements.
	Proxy      string `yaml:"proxy,omitempty"`
	HTTPProxy  string `yaml:"http_proxy,omitempty"`
	FtpProxy   string `yaml:"ftp_proxy,omitempty"`
	HTTPSProxy string `yaml:"https_proxy,omitempty"`
	// 'source' entries in apt-sources that match this python regex
	// expression will be passed to add-apt-repository
	// default:  '^[\w-]+:\w'
	AddAptRepoMatch string `yaml:"add_apt_repo_match,omitempty"`
	// The key of each source entry is the filename and will be prepended by
	// /etc/apt/sources.list.d/ if it doesn't start with a '/'.
	// If it doesn't end with .list it will be appended so that apt picks up it's
	// configuration.
	Sources map[string]AptSource `yaml:"sources,omitempty"`
}

type File struct {
	//b64 or gzip or (gz+b64)
	Encoding    string      `yaml:"encoding,omitempty"`
	Content     interface{} `yaml:"content,omitempty"`
	Owner       string      `yaml:"owner,omitempty"`
	Path        string      `yaml:"path,omitempty"`
	Permissions string      `yaml:"permissions,omitempty"`
}

type FsSetup struct {
	// The file system label to be used. If set to None, no label is  used.
	Label string `yaml:"label,omitempty"`
	// The file system type. It is assumed that the there
	//        will be a "mkfs.<FS_TYPE>" that behaves likes "mkfs". On a standard
	//        Ubuntu Cloud Image, this means that you have the option of ext{2,3,4},and vfat by default.
	//
	Filesystem string `yaml:"filesystem,omitempty"`

	//The device name. Special names of 'ephemeralX' or 'swap'
	//        are allowed and the actual device is acquired from the cloud datasource.
	//        When using 'ephemeralX' (i.e. ephemeral0), make sure to leave the
	//        label as 'ephemeralX' otherwise there may be issues with the mounting
	//        of the ephemeral storage layer.
	//
	//        If you define the device as 'ephemeralX.Y' then Y will be interpetted
	//        as a partition value. However, ephermalX.0 is the _same_ as ephemeralX.
	Device string `yaml:"device,omitempty"`
	// Partition definitions are overwritten if you use the '<DEVICE>.Y' notation.
	//
	//        The valid options are:
	//        "auto|any": tell cloud-init not to care whether there is a partition
	//            or not. Auto will use the first partition that does not contain a
	//            file system already. In the absence of a partition table, it will
	//            put it directly on the disk.
	//
	//            "auto": If a file system that matches the specification in terms of
	//            label, type and device, then cloud-init will skip the creation of
	//            the file system.
	//
	//            "any": If a file system that matches the file system type and device,
	//            then cloud-init will skip the creation of the file system.
	//
	//            Devices are selected based on first-detected, starting with partitions
	//            and then the raw disk. Consider the following:
	//                NAME     FSTYPE LABEL
	//                xvdb
	//                |-xvdb1  ext4
	//                |-xvdb2
	//                |-xvdb3  btrfs  test
	//                \-xvdb4  ext4   test
	//
	//            If you ask for 'auto', label of 'test, and file system of 'ext4'
	//            then cloud-init will select the 2nd partition, even though there
	//            is a partition match at the 4th partition.
	//
	//            If you ask for 'any' and a label of 'test', then cloud-init will
	//            select the 1st partition.
	//
	//            If you ask for 'auto' and don't define label, then cloud-init will
	//            select the 1st partition.
	//
	//            In general, if you have a specific partition configuration in mind,
	//            you should define either the device or the partition number. 'auto'
	//            and 'any' are specifically intended for formating ephemeral storage or
	//            for simple schemes.
	//
	//        "none": Put the file system directly on the device.
	//
	//        <NUM>: where NUM is the actual partition number.
	Partition string `yaml:"partition,omitempty"`
	//	Defines whether or not to overwrite any existing
	//        filesystem.
	//
	//        "true": Indiscriminately destroy any pre-existing file system. Use at
	//            your own peril.
	//
	//        "false": If an existing file system exists, skip the creation.
	Overwrite bool `yaml:"overwrite,omitempty"`
	//This is a special directive, used for Windows Azure that
	//        instructs cloud-init to replace a file system of <FS_TYPE>. NOTE:
	//        unless you define a label, this requires the use of the 'any' partition
	//        directive.
	ReplaseFS bool   `yaml:"replace_fs,omitempty"`
	Cmd       string `yaml:"cmd,omitempty"`
}

type Chef struct {
	InstallType       string   `yaml:"install_type,omitempty"`
	ForceInstall      bool     `yaml:"force_install,omitempty"`
	ServerURL         string   `yaml:"server_url,omitempty"`
	NodeName          string   `yaml:"node_name,omitempty"`
	Environment       string   `yaml:"environment,omitempty"`
	ValidationName    string   `yaml:"validation_name,omitempty"`
	ValidationCert    string   `yaml:"validation_cert,omitempty"`
	RunList           []string `yaml:"run_list,omitempty"`
	InitialAttributes struct {
		Apache struct {
			Prefork struct {
				Maxclients int `yaml:"maxclients,omitempty"`
			} `yaml:"prefork,omitempty"`
			Keepalive string `yaml:"keepalive,omitempty"`
		} `yaml:"apache,omitempty"`
	} `yaml:"initial_attributes,omitempty"`
	OmnibusURL     string `yaml:"omnibus_url,omitempty"`
	OmnibusVersion string `yaml:"omnibus_version,omitempty"`
}

// power_state can be used to make the system shutdown, reboot or
// halt after boot is finished.  This same thing can be achieved by
// user-data scripts or by runcmd by simply invoking 'shutdown'.
//
// Doing it this way ensures that cloud-init is entirely finished with
// modules that would be executed, and avoids any error/log messages
// that may go to the console as a result of system services like
// syslog being taken down while cloud-init is running.
//
// If you delay '+5' (5 minutes) and have a timeout of
// 120 (2 minutes), then the max time until shutdown will be 7 minutes.
// cloud-init will invoke 'shutdown +5' after the process finishes, or
// when 'timeout' seconds have elapsed.
//
type PowerState struct {
	//form accepted by shutdown.  default is 'now'. other format  accepted is +m (m in minutes)
	Delay string `yaml:"delay,omitempty"`
	//required. must be one of 'poweroff', 'halt', 'reboot'
	Mode string `yaml:"mode,omitempty"`
	//provided as the message argument to 'shutdown'. default is none.
	Message string `yaml:"message,omitempty"`
	//the amount of time to give the cloud-init process to finish  before executing shutdown.
	Timeout int `yaml:"timeout,omitempty"`

	// apply state change only if condition is met.
	//    May be boolean True (always met), or False (never met),
	//    or a command string or list to be executed.
	//    command's exit code indicates:
	//       0: condition met
	//       1: condition not met
	//    other exit codes will result in 'not met', but are reserved  for future use.
	Condition string `yaml:"condition,omitempty"`
}

type MAASDatasource struct {
	Timeout     int    `yaml:"timeout,omitempty"`
	MaxWait     int    `yaml:"max_wait,omitempty"`
	MetadataURL string `yaml:"metadata_url,omitempty"`
	ConsumerKey string `yaml:"consumer_key,omitempty"`
	TokenKey    string `yaml:"token_key,omitempty"`
	TokenSecret string `yaml:"token_secret,omitempty"`
}

type Ec2Datasource struct {
	Timeout      int      `yaml:"timeout,omitempty"`
	MaxWait      int      `yaml:"max_wait,omitempty"`
	MetadataUrls []string `yaml:"metadata_urls,omitempty"`
}

type NoCloudMetadata struct {
	InstanceID    string `yaml:"instance-id,omitempty"`
	LocalHostname string `yaml:"local-hostname,omitempty"`
}
type NoCloudDatasource struct {
	Seedfrom string          `yaml:"seedfrom,omitempty"`
	FsLabel  string          `yaml:"fs_label,omitempty"`
	UserData string          `yaml:"user-data,omitempty"`
	MetaData NoCloudMetadata `yaml:"meta-data,omitempty"`
}

type AzureDatasource struct {
	AgentCommand   []string `yaml:"agent_command,omitempty"`
	SetHostname    bool     `yaml:"set_hostname,omitempty"`
	HostnameBounce struct {
		Interface string `yaml:"interface,omitempty"`
		Policy    bool   `yaml:"policy,omitempty"`
	} `yaml:"hostname_bounce,omitempty"`
}

type SmartOSDatasource struct {
	SerialDevice     string        `yaml:"serial_device,omitempty"`
	SerialTimeout    int           `yaml:"serial_timeout,omitempty"`
	MetadataSockfile string        `yaml:"metadata_sockfile,omitempty"`
	NoBase64Decode   []string      `yaml:"no_base64_decode,omitempty"`
	Base64Keys       []interface{} `yaml:"base64_keys,omitempty"`
	Base64All        bool          `yaml:"base64_all,omitempty"`
}
type Datasource struct {
	Ec2     Ec2Datasource     `yaml:"Ec2,omitempty"`
	MAAS    MAASDatasource    `yaml:"MAAS,omitempty"`
	NoCloud NoCloudDatasource `yaml:"NoCloud,omitempty"`
	Azure   AzureDatasource   `yaml:"Azure,omitempty"`
	SmartOS SmartOSDatasource `yaml:"SmartOS,omitempty"`
}

type PuppetAgent struct {
	Server   string `yaml:"server,omitempty"`
	Certname string `yaml:"certname,omitempty"`
}

type PuppetConf struct {
	Agent  PuppetAgent `yaml:"agent,omitempty"`
	CaCert string      `yaml:"ca_cert,omitempty"`
}
type Puppet struct {
	Conf PuppetConf `yaml:"conf,omitempty"`
}

type Phonehome struct {
	URL string `yaml:"url,omitempty"`
	// [ pub_key_dsa, pub_key_rsa, pub_key_ecdsa, instance_id ]
	Post  []string `yaml:"post,omitempty"`
	Tries int      `yaml:"tries,omitempty"`
}

type Swap struct {
	Filename string `yaml:"filename,omitempty"`
	Size     string `yaml:"size,omitempty"`
	Maxsize  string `yaml:"maxsize,omitempty"`
}
type User struct {
	// The user's login name
	Name string `yaml:"name,omitempty"`
	//The user name's real name, i.e. "Bob B. Smith"
	Gecos string `yaml:"gecos,omitempty"`
	//	Optional. The SELinux user for the user's login, such as
	//          "staff_u". When this is omitted the system will select the default
	//           SELinux user.
	SeLinuxUser string `yaml:"selinux_user,omitempty"`
	ExpireDate  string `yaml:"expiredate,omitempty"`
	//	Defaults to none. Accepts a sudo rule string, a list of sudo rule
	//         strings or False to explicitly deny sudo usage. Examples:
	//
	//         Allow a user unrestricted sudo access.
	//             sudo:  ALL=(ALL) NOPASSWD:ALL
	//
	//         Adding multiple sudo rule strings.
	//             sudo:
	//               - ALL=(ALL) NOPASSWD:/bin/mysql
	//               - ALL=(ALL) ALL
	//
	//         Prevent sudo access for a user.
	//             sudo: False
	//
	//         Note: Please double check your syntax and make sure it is valid.
	//               cloud-init does not parse/check the syntax of the sudo
	//               directive.
	Sudo string `yaml:"sudo,omitempty"`
	//	The hash -- not the password itself -- of the password you want
	//           to use for this user. You can generate a safe hash via:
	//               mkpasswd --method=SHA-512 --rounds=4096
	//           (the above command would create from stdin an SHA-512 password hash
	//           with 4096 salt rounds)
	//
	//           Please note: while the use of a hashed password is better than
	//               plain text, the use of this feature is not ideal. Also,
	//               using a high number of salting rounds will help, but it should
	//               not be relied upon.
	//
	//               To highlight this risk, running John the Ripper against the
	//               example hash above, with a readily available wordlist, revealed
	//               the true password in 12 seconds on a i7-2620QM.
	//
	//               In other words, this feature is a potential security risk and is
	//               provided for your convenience only. If you do not fully trust the
	//               medium over which your cloud-config will be transmitted, then you
	//               should use SSH authentication only.
	//
	//               You have thus been warned.
	Passwd string `yaml:"passwd,omitempty"`
	// define the primary group. Defaults to a new group created named after the user.
	PrimaryGroup string `yaml:"primary_group,omitempty"`
	Groups       string `yaml:"groups,omitempty"`
	// Optional. Import SSH ids
	SSHImportID string `yaml:"ssh_import_id,omitempty"`
	//Defaults to true. Lock the password to disable password login
	LockPasswd *bool `yaml:"lock_passwd,omitempty"`
	//When set to true, do not create home directory
	NoCreateHome bool `yaml:"no_create_home,omitempty"`
	//When set to true, do not create a group named after the user.
	NoUserGroup bool `yaml:"no_user_group,omitempty"`
	//When set to true, do not initialize lastlog and faillog database.
	NoLogInit bool `yaml:"no_log_init,omitempty"`
	//Add keys to user's authorized keys file
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	//Create the user as inactive
	Inactive bool `yaml:"inactive,omitempty"`
	// Create the user as a system user. This means no home directory.
	System bool `yaml:"system,omitempty"`
	//Create a Snappy (Ubuntu-Core) user via the snap create-user
	//             command available on Ubuntu systems.  If the user has an account
	//             on the Ubuntu SSO, specifying the email will allow snap to
	//             request a username and any public ssh keys and will import
	//             these into the system with username specified by SSO account./
	//             If 'username' is not set in SSO, then username will be the
	//             shortname before the email domain.
	Snapuser string `yaml:"snapuser,omitempty"`
	//	Set true to block ssh logins for cloud
	//      ssh public keys and emit a message redirecting logins to
	//      use <default_username> instead. This option only disables cloud
	//      provided public-keys. An error will be raised if ssh_authorized_keys
	//      or ssh_import_id is provided for the same user.
	SSHRedirectUser bool `yaml:"ssh_redirect_user,omitempty"`
}

type YumRepo struct {
	Baseurl        string `yaml:"baseurl,omitempty" validate:"required"`
	Enabled        bool   `yaml:"enabled,omitempty"`
	Failovermethod string `yaml:"failovermethod,omitempty"`
	Gpgcheck       bool   `yaml:"gpgcheck,omitempty"`
	Gpgkey         string `yaml:"gpgkey,omitempty"`
	Name           string `yaml:"name,omitempty"`
}

type SSHKeys struct {
	RsaPrivate   string `yaml:"rsa_private,omitempty"`
	RsaPublic    string `yaml:"rsa_public,omitempty"`
	DsaPrivate   string `yaml:"dsa_private,omitempty"`
	DsaPublic    string `yaml:"dsa_public,omitempty"`
	EcdsaPublic  string `yaml:"ecdsa_public,omitempty"`
	EcdsaPrivate string `yaml:"ecdsa_private,omitempty"`
}
type RhSubscription struct {
	Username      string `yaml:"username,omitempty"`
	Password      string `yaml:"password,omitempty"`
	ActivationKey string `yaml:"activation-key,omitempty"`
	Org           int    `yaml:"org,omitempty"`
	AutoAttach    bool   `yaml:"auto-attach,omitempty"`
	//to set the service level for your   subscriptions
	ServiceLevel string `yaml:"service-level,omitempty"`
	//	needs to be a list of repo IDs
	EnableRepo []interface{} `yaml:"enable-repo,omitempty"`
	//	needs to be a list of repo IDs
	DisableRepo []interface{} `yaml:"disable-repo,omitempty"`
	//alter baseurl in /etc/rhsm/rhsm.conf
	RhsmBaseurl string `yaml:"rhsm-baseurl,omitempty"`
	//alter the server hostname in   /etc/rhsm/rhsm.conf
	ServerHostname string `yaml:"server-hostname,omitempty"`
}

type Growpart struct {
	Mode                   string   `yaml:"mode,omitempty"`
	Devices                []string `yaml:"devices,omitempty"`
	IgnoreGrowrootDisabled bool     `yaml:"ignore_growroot_disabled,omitempty"`
}

type DiskSetup struct {
	TableType string `yaml:"table_type,omitempty"`
	Layout    bool   `yaml:"layout,omitempty"`
	Overwrite bool   `yaml:"overwrite,omitempty"`
}

type CloudInit struct {
	FileEncoding string    `yaml:"-"`
	FsSetup      []FsSetup `yaml:"fs_setup,omitempty"`
	Chef         Chef      `yaml:"chef,omitempty"`
	Output       struct {
		All string `yaml:"all,omitempty"`
	} `yaml:"output,omitempty"`
	DeviceAliases map[string]string `yaml:"device_aliases,omitempty"`
	PowerState    PowerState        `yaml:"power_state,omitempty"`
	//	 default: cloud-init boot finished at $TIMESTAMP. Up $UPTIME seconds
	// this message is written by cloud-final when the system is finished
	// its first boot
	FinalMessage string `yaml:"final_message,omitempty"`
	// 'mounts' contains a list of lists
	//  the inner list are entries for an /etc/fstab line
	//  ie : [ fs_spec, fs_file, fs_vfstype, fs_mntops, fs-freq, fs_passno ]
	//
	// default:
	// mounts:
	//  - [ ephemeral0, /mnt ]
	//  - [ swap, none, swap, sw, 0, 0 ]
	//
	// in order to remove a previously listed mount (ie, one from defaults)
	// list only the fs_spec.  For example, to override the default, of
	// mounting swap:
	// - [ swap ]
	// or
	// - [ swap, null ]
	//
	// - if a device does not exist at the time, an entry will still be
	//   written to /etc/fstab.
	// - '/dev' can be omitted for device names that begin with: xvd, sd, hd, vd
	// - if an entry does not have all 6 fields, they will be filled in
	//   with values from 'mount_default_fields' below.
	//
	// Note, that you should set 'nofail' (see man fstab) for volumes that may not
	// be attached at instance boot (or reboot).
	Mounts [][]string `yaml:"mounts,omitempty"`
	Puppet Puppet     `yaml:"puppet,omitempty"`
	// These values are used to fill in any entries in 'mounts' that are not
	// complete.  This must be an array, and must have 7 fields.
	MountDefaultFields []string `yaml:"mount_default_fields,omitempty"`
	// swap can also be set up by the 'mounts' module
	// default is to not create any swap files, because 'size' is set to 0
	Swap           Swap                 `yaml:"swap,omitempty"`
	RhSubscription RhSubscription       `yaml:"rh_subscription,omitempty"`
	DiskSetup      map[string]DiskSetup `yaml:"disk_setup,omitempty"`
	Growpart       Growpart             `yaml:"growpart,omitempty"`
	Phonehome      Phonehome            `yaml:"phone_home,omitempty"`
	Datasource     Datasource           `yaml:"datasource,omitempty"`
	Groups         []interface{}        `yaml:"groups,omitempty"`
	//this is very similar to runcmd, but commands run very early
	// in the boot process, only slightly after a 'boothook' would run.
	// bootcmd should really only be used for things that could not be
	// done later in the boot process.  bootcmd is very much like
	// boothook, but possibly with more friendly.
	// - bootcmd will run on every boot
	// - the INSTANCE_ID variable will be set to the current instance id.
	// - you can use 'cloud-init-per' command to help only run once
	Bootcmd []interface{} `yaml:"bootcmd,omitempty"`
	// - runcmd only runs during the first boot
	// - if the item is a list, the items will be properly executed as if
	//   passed to execve(3) (with the first arg as the command).
	// - if the item is a string, it will be simply written to the file and
	//   will be interpreted by 'sh'
	//
	// Note, that the list has to be proper yaml, so you have to quote
	// any characters yaml would eat (':' can be problematic)
	Runcmd [][]string `yaml:"runcmd,omitempty"`
	// if packages are specified, this apt_update will be set to true
	// packages may be supplied as a single package name or as a list
	// with the format [<package>, <version>] wherein the specific package version will be installed.
	Packages []interface{} `yaml:"packages,omitempty"`
	// add each entry to ~/.ssh/authorized_keys for the configured user or th  first user defined in the user definition directive.
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	// Send pre-generated ssh private keys to the server
	// If these are present, they will be written to /etc/ssh and new random keys will not be generated
	SSHKeys SSHKeys `yaml:"ssh_keys,omitempty"`
	// Update apt database on first boot (run 'apt-get update').
	// Note, if packages are given, or package_upgrade is true, then
	// update will be done independent of this setting.
	PackageUpdate bool `yaml:"package_update,omitempty"`
	// Upgrade the instance on first boot (ie run apt-get upgrade)
	PackageUpgrade   bool   `yaml:"package_upgrade,omitempty"`
	PreserveHostname bool   `yaml:"preserve_hostname,omitempty"`
	Hostname         string `yaml:"hostname,omitempty"`
	Users            []User `yaml:"users,omitempty"`
	WriteFiles       []File `yaml:"write_files,omitempty"`
	Apt              Apt    `yaml:"apt,omitempty"`
	// Configure Acquire::http::Pipeline-Depth
	// Default: disables HTTP pipelining. Certain web servers, such as S3 do not pipeline properly (LP: #948461).
	// Valid options:
	//   False/default: Disables pipelining for APT
	//   None/Unchanged: Use OS default
	//   Number: Set pipelining to some number (not recommended)
	AptPipelining string             `yaml:"apt_pipelining,omitempty"`
	YumRepos      map[string]YumRepo `yaml:"yum_repos,omitempty"`
}
