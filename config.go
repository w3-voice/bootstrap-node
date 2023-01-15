package hoodboot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/hood-chat/core"
	"github.com/libp2p/go-libp2p/core/peer"

	ma "github.com/multiformats/go-multiaddr"
)

// Config stores the full configuration of the relays, ACLs and other settings
// that influence behaviour of a relay daemon.
type Config struct {
	Network   NetworkConfig
	ConnMgr   ConnMgrConfig
	Bootstrap BootstrapConfig
	Daemon    DaemonConfig
}

// DaemonConfig controls settings for the relay-daemon itself.
type DaemonConfig struct {
	PprofPort int
}

// NetworkConfig controls listen and annouce settings for the libp2p host.
type NetworkConfig struct {
	ListenAddrs   []string
	AnnounceAddrs []string
}

// ConnMgrConfig controls the libp2p connection manager settings.
type ConnMgrConfig struct {
	ConnMgrLo    int
	ConnMgrHi    int
	ConnMgrGrace time.Duration
}


// Bootstrap Config
type BootstrapConfig struct {
	Peers             []string
	MinPeerThreshold  int
}



// DefaultConfig returns a default relay configuration using default resource
// settings and no ACLs.
func DefaultConfig() Config {
	return Config{
		Network: NetworkConfig{
			ListenAddrs: []string{
				"/ip4/0.0.0.0/udp/4002/quic",
				"/ip6/::/udp/4001/quic",
				"/ip4/0.0.0.0/tcp/4002",
				"/ip6/::/tcp/4001",
			},
		},
		ConnMgr: ConnMgrConfig{
			ConnMgrLo:    512,
			ConnMgrHi:    768,
			ConnMgrGrace: 2 * time.Minute,
		},
		Bootstrap:BootstrapConfig{
			Peers: core.BootstrapNodes,
			MinPeerThreshold: 1,
		},
		Daemon: DaemonConfig{
			PprofPort: 6060,
		},
	}
}

// LoadConfig reads a relay daemon JSON configuration from the given path.
// The configuration is first initialized with DefaultConfig, so all unset
// fields will take defaults from there.
func LoadConfig(cfgPath string) (Config, error) {
	cfg := DefaultConfig()

	if cfgPath != "" {
		cfgFile, err := os.Open(cfgPath)
		if err != nil {
			return Config{}, err
		}
		defer cfgFile.Close()

		decoder := json.NewDecoder(cfgFile)
		err = decoder.Decode(&cfg)
		if err != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}

func ParsePeers(addrs []string) ([]peer.AddrInfo, error) {
	maddrs := make([]ma.Multiaddr, len(addrs))
	for i, addr := range addrs {
		var err error
		maddrs[i], err = ma.NewMultiaddr(addr)
		if err != nil {
			return nil, err
		}
	}
	return peer.AddrInfosFromP2pAddrs(maddrs...)
}
