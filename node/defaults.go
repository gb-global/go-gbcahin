package node

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"gbchain-org/go-gbchain/p2p"
	"gbchain-org/go-gbchain/p2p/nat"
	"gbchain-org/go-gbchain/rpc"
)

const (
	DefaultHTTPHost    = "localhost" // Default host interface for the HTTP RPC server
	DefaultHTTPPort    = 8585        // Default TCP port for the HTTP RPC server
	DefaultWSHost      = "localhost" // Default host interface for the websocket RPC server
	DefaultWSPort      = 8586        // Default TCP port for the websocket RPC server
	DefaultGraphQLHost = "localhost" // Default host interface for the GraphQL server
	DefaultGraphQLPort = 8587        // Default TCP port for the GraphQL server
	DefaultSubHTTPHost = "localhost" // Default host interface for the HTTP RPC server
	DefaultSubHTTPPort = 9585        // Default TCP port for the HTTP RPC server
	DefaultSubWSHost   = "localhost" // Default host interface for the websocket RPC server
	DefaultSubWSPort   = 9586        // Default TCP port for the websocket RPC server
)

// DefaultConfig contains reasonable default settings.
var DefaultConfig = Config{
	DataDir:             DefaultDataDir(),
	HTTPPort:            DefaultHTTPPort,
	HTTPModules:         []string{"net", "web3"},
	HTTPVirtualHosts:    []string{"localhost"},
	HTTPTimeouts:        rpc.DefaultHTTPTimeouts,
	WSPort:              DefaultWSPort,
	WSModules:           []string{"net", "web3"},
	GraphQLPort:         DefaultGraphQLPort,
	GraphQLVirtualHosts: []string{"localhost"},
	P2P: p2p.Config{
		ListenAddr: ":36608",
		MaxPeers:   30,
		NAT:        nat.Any(),
	},
	SubHTTPPort:         DefaultSubHTTPPort,
	SubHTTPModules:      []string{"net", "web3"},
	SubHTTPVirtualHosts: []string{"localhost"},
	SubWSPort:           DefaultSubWSPort,
	SubWSModules:        []string{"net", "web3"},
}

// DefaultDataDir is the default data directory to use for the databases and other
// persistence requirements.
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		switch runtime.GOOS {
		case "darwin":
			return filepath.Join(home, "Library", "GBChain")
		case "windows":
			// We used to put everything in %HOME%\AppData\Roaming, but this caused
			// problems with non-typical setups. If this fallback location exists and
			// is non-empty, use it, otherwise DTRT and check %LOCALAPPDATA%.
			fallback := filepath.Join(home, "AppData", "Roaming", "GBChain")
			appdata := windowsAppData()
			if appdata == "" || isNonEmptyDir(fallback) {
				return fallback
			}
			return filepath.Join(appdata, "GBChain")
		default:
			return filepath.Join(home, ".GBChain")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

func windowsAppData() string {
	v := os.Getenv("LOCALAPPDATA")
	if v == "" {
		// Windows XP and below don't have LocalAppData. Crash here because
		// we don't support Windows XP and undefining the variable will cause
		// other issues.
		panic("environment variable LocalAppData is undefined")
	}
	return v
}

func isNonEmptyDir(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	names, _ := f.Readdir(1)
	f.Close()
	return len(names) > 0
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}
