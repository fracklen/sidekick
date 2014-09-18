package sidekick

// ContainerInfo - all the information needed about a container
type ContainerInfo struct {
	ID      string
	Created string
	Path    string
	Name    string
	Args    []string
	Config  *containerConfig
	State   struct {
		Running   bool
		Pid       int
		ExitCode  int
		StartedAt string
		Ghost     bool
	}
	Image           string
	NetworkSettings struct {
		IPAddress   string
		IPPrefixLen int
		Gateway     string
		Bridge      string
		Ports       map[string][]portBinding
	}
	SysInitPath    string
	ResolvConfPath string
	Volumes        map[string]string
	HostConfig     *hostConfig
}

type containerConfig struct {
	Hostname        string
	Domainname      string
	User            string
	Memory          int
	MemorySwap      int
	CPUShares       int
	CPUset          string
	AttachStdin     bool
	AttachStdout    bool
	AttachStderr    bool
	PortSpecs       []string
	ExposedPorts    map[string]struct{}
	Tty             bool
	OpenStdin       bool
	StdinOnce       bool
	Env             []string
	Cmd             []string
	Image           string
	Volumes         map[string]struct{}
	WorkingDir      string
	Entrypoint      []string
	NetworkDisabled bool
	OnBuild         []string
}

type hostConfig struct {
	Binds           []string
	ContainerIDFile string
	LxcConf         []map[string]string
	Privileged      bool
	PortBindings    map[string][]portBinding
	Links           []string
	PublishAllPorts bool
	DNS             []string
	DNSSearch       []string
	VolumesFrom     []string
	NetworkMode     string
}

type portBinding struct {
	HostIP   string
	HostPort string
}

type port struct {
	PrivatePort int
	PublicPort  int
	Type        string
}

type container struct {
	ID         string
	Names      []string
	Image      string
	Command    string
	Created    int
	Status     string
	Ports      []port
	SizeRw     int
	SizeRootFs int
}

type event struct {
	ID     string
	Status string
	From   string
	Time   int
}

type version struct {
	Version   string
	GitCommit string
	GoVersion string
}

type respContainersCreate struct {
	ID       string
	Warnings []string
}
