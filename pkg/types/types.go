package types

import (
	"fmt"
	"strings"

	"github.com/moshloop/konfigadm/pkg/os"

	cloudinit "github.com/moshloop/konfigadm/pkg/cloud-init"
	"github.com/moshloop/konfigadm/pkg/utils"
)

var (
	konfigadm = "konfigadm"
)

//Port maps src and target ports
type Port struct {
	Port   int `json:"port,omitempty"  validate:"min=1,max=65536"`
	Target int `json:"target,omitempty"  validate:"min=1,max=65536"`
}

//Container represents a container to be run using systemd
type Container struct {
	//The name of the service (e.g systemd unit name or deployment name)
	Service string `json:"service,omitempty"`

	Image string `json:"image"`

	//A map of environment variables to pass through
	Env map[string]string `json:"env,omitempty"`

	//A map of labels to add to the container
	Labels map[string]string `json:"labels,omitempty"`

	//Additional arguments to the docker run command e.g. -p 8080:8080
	DockerOpts string `json:"docker_opts,omitempty"`

	//Additional options to the docker client e.g. -H unix:///tmp/var/run/docker.sock
	DockerClientArgs string `json:"docker_client_args,omitempty"`

	//Additional arguments to the container
	Args string `json:"args,omitempty"`

	Ports []Port `json:"ports,omitempty"`

	Commands []string `json:"commands,omitempty"`

	//Map of files to mount into the container
	Files map[string]string `json:"files,omitempty"`

	//Map of templates to mount into the container
	Templates map[string]string `json:"templates,omitempty"`

	//TODO:
	Volumes []string `json:"volumes,omitempty"`

	//TODO  capabilities:

	//CPU limit in cores (Defaults to 1 )
	CPU int `json:"cpu,omitempty" validate:"min=0,max=32"`

	//	Memory Limit in MB. (Defaults to 1024)
	Mem int `json:"mem,omitempty" validate:"min=0,max=1048576"`

	//default:	user-bridge	 only
	Network string `json:"network,omitempty"`

	// default: 1
	Replicas int `json:"replicas,omitempty"`
}

func (c Container) Name() string {
	if c.Service != "" {
		return c.Service
	}
	name := strings.Split(c.Image, ":")[0]
	if strings.Contains(name, "/") {
		name = name[strings.LastIndex(name, "/")+1:]
	}
	return name
}


//User mirrors the CloudInit User struct.
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
	LockPasswd bool `yaml:"lock_passwd,omitempty"`
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

//File is a primitive representing a single file
type File struct {
	Content        string   `yaml:"content,omitempty"`
	ContentFromURL string   `yaml:"content_from_url,omitempty"`
	Unarchive      bool     `yaml:"unarchive,omitempty"`
	Permissions    string   `yaml:"permissions,omitempty"`
	Owner          string   `yaml:"owner,omitempty"`
	Flags          []Flag `yaml:"tags,omitempty"`
}

//Filesystem is a primitive for referencing all files
type Filesystem map[string]File

type Certificate string

//Config is the logical model after runtime tags have been applied
type Config struct {

	/** Primitive elements are what all native and operator commands eventually compile down into **/

	PreCommands  []Command  `yaml:"pre_commands,omitempty"`
	Commands     []Command  `yaml:"commands,omitempty"`
	PostCommands []Command  `yaml:"post_commands,omitempty"`
	Filesystem   Filesystem `yaml:"filesystem,omitempty"`

	/** Native elements are "compiled" into primitive items in order to apply them **/

	//Files is a map of destination path to lookup file path
	// The lookup path is relative to where konfigadm is run from, not relative to the config file
	// The content and permissions of the file will be compiled into primitive Filesystem elements, user and group ownership is ignored
	// Both the destination and lookup path can be expressions
	Files map[string]string `yaml:"files,omitempty"`

	//Templates is a map of destination path to template lookup path
	// The lookup path is relative to where konfigadm is run from, not relative to the config file
	// Templates are compiled via a Jinja (Ansible-like) rendered into primitive filesystem objects
	// Both the destination and lookup path can be expressions
	Templates        map[string]string    `yaml:"templates,omitempty"`
	Sysctls          map[string]string    `yaml:"sysctls,omitempty"`
	Packages         *[]Package           `yaml:"packages,omitempty"`
	PackageRepos     *[]PackageRepo       `yaml:"package_repos,omitempty"`
	Images           []string             `yaml:"images,omitempty"`
	Containers       []Container          `yaml:"containers,omitempty"`
	ContainerRuntime *ContainerRuntime    `yaml:"container_runtime,omitempty"`
	Kubernetes       *KubernetesSpec      `yaml:"kubernetes,omitempty"`
	Environment      map[string]string    `yaml:"environment,omitempty"`
	Timezone         string               `yaml:"timezone,omitempty"`
	NTP              []string             `yaml:"ntp,omitempty"`
	DNS              []string             `yaml:"dns,omitempty"`
	Limits           []string             `yaml:"limits,omitempty"`
	TrustedCA        []Certificate        `yaml:"ca,omitempty"`
	Partitions       []string             `yaml:"partitions,omitempty"`
	Extra            *cloudinit.CloudInit `yaml:"extra,omitempty"`
	Services         map[string]Service   `yaml:"services,omitempty"`
	Users            []User               `yaml:"users,omitempty"`
	Context          *SystemContext       `yaml:"-"`
}

type Applier interface {
	Apply(ctx SystemContext)
}

type SystemContext struct {
	Vars  map[string]interface{}
	Flags []Flag
	Name  string
	OS    os.OS
}

type Transformer func(cfg *Config, ctx *SystemContext) (commands []Command, files Filesystem, err error)

type FlagProcessor func(cfg *Config, flags ...Flag)

type AllPhases interface {
	Phase
	ProcessFlagsPhase
}

type Phase interface {
	ApplyPhase(cfg *Config, ctx *SystemContext) (commands []Command, files Filesystem, err error)
}

type ProcessFlagsPhase interface {
	ProcessFlags(cfg *Config, flags ...Flag)
}

type VerifyPhase interface {
	Verify(cfg *Config, results *VerifyResults, flags ...Flag) bool
}

//Results records the results of a test or verification run
type VerifyResults struct {
	PassCount int
	FailCount int
	SkipCount int
}

func (c *VerifyResults) Done() {
	fmt.Printf("  %d passed, %d skipped, %d failed\n", c.PassCount, c.SkipCount, c.FailCount)
}

func (c *VerifyResults) Pass(msg string, args ...interface{}) {
	c.PassCount++
	fmt.Printf("%s [pass] %s %s\n", utils.Green, fmt.Sprintf(msg, args...), utils.Reset)

}
func (c *VerifyResults) Fail(msg string, args ...interface{}) {
	c.FailCount++
	fmt.Printf("%s [fail] %s %s\n", utils.Red, fmt.Sprintf(msg, args...), utils.Reset)
}
func (c VerifyResults) Skip(msg string, args ...interface{}) {
	c.SkipCount++
	fmt.Printf("%s [skip] %s %s\n", utils.LightCyan, fmt.Sprintf(msg, args...), utils.Reset)
}
