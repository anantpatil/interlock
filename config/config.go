package config

// ExtensionConfig has all options for all load balancer extensions
// the extension itself will use whichever options needed
type ExtensionConfig struct {
	Name                   string // extension name
	ConfigPath             string // config file path
	ConfigBasePath         string `toml:"-"` // internal
	PidPath                string // haproxy, nginx
	TemplatePath           string // template file path
	BackendOverrideAddress string // haproxy, nginx
	ConnectTimeout         int    // haproxy
	ServerTimeout          int    // haproxy
	ClientTimeout          int    // haproxy
	MaxConn                int    // haproxy, nginx
	Port                   int    // haproxy, nginx
	SyslogAddr             string // haproxy
	AdminUser              string // haproxy
	AdminPass              string // haproxy
	SSLCertPath            string // haproxy, nginx
	SSLCert                string // haproxy
	SSLPort                int    // haproxy, nginx
	SSLOpts                string // haproxy
	SSLDefaultDHParam      int    // haproxy
	SSLServerVerify        string // haproxy
	DHParam                bool   // nginx
	DHParamPath            string // nginx
	NginxPlusEnabled       bool   // nginx
	User                   string // nginx
	WorkerProcesses        int    // nginx
	RLimitNoFile           int    // nginx
	ProxyConnectTimeout    int    // nginx
	ProxySendTimeout       int    // nginx
	ProxyReadTimeout       int    // nginx
	SendTimeout            int    // nginx
	SSLCiphers             string // nginx
	SSLProtocols           string // nginx
}

// Config is the top level configuration
type Config struct {
	ListenAddr    string
	GRPCAddr      string
	DockerURL     string
	TLSCACert     string
	TLSCert       string
	TLSKey        string
	AllowInsecure bool
	EnableMetrics bool
	PollInterval  string
	Extensions    []*ExtensionConfig
}
