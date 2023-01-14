package main

import (
	"github.com/hood-chat/core"
	"github.com/hood-chat/core/entity"
	"github.com/hood-chat/core/repo"
	"github.com/hood-chat/core/store"
	libp2p "github.com/libp2p/go-libp2p"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p/config"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/sys/unix"
)

var log = logging.Logger("boothood")

// Main function
func main() {
	s, err := store.NewStore("./data")
	if err != nil {
		panic(err)
	}
	rIdentity := repo.NewIdentityRepo(s)
	id, err := rIdentity.Get()
	if err != nil {
		id, err = entity.CreateIdentity("hoodboot")
		if err != nil {
			panic(err)
		}
		err := rIdentity.Set(id)
		if err != nil {
			panic(err)
		}

	}
	opt := Option()

	err = opt.SetIdentity(&id)
	if err != nil {
		log.Debugf("Can not store identity")
		panic("can not store identity")
	}
	hb := core.DefaultRoutedHost{}
	if err != nil {
		panic(err)
	}
	h, err := hb.Create(opt)

	if err != nil {
		panic(err)
	}

	log.Warn("Hoodboot listens on %s", h.Addrs())
	log.Warn("Hoodboot Peer ID is %s", h.ID())

	select {} // block forever
}

var ListenAddrs = func(cfg *config.Config) error {
	ip4ListenAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	if err != nil {
		return err
	}
	quicListenAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/udp/4001/quic")
	if err != nil {
		return err
	}
	quicV1ListenAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/udp/4001/quic-v1")
	if err != nil {
		return err
	}
	defaultIP6ListenAddr, err := multiaddr.NewMultiaddr("/ip6/::/tcp/4001")
	if err != nil {
		return err
	}

	return cfg.Apply(libp2p.ListenAddrs(
		quicListenAddr,
		ip4ListenAddr,
		quicV1ListenAddr,
		defaultIP6ListenAddr,
	))
}

func Option() core.Option {

	opt := []libp2p.Option{
		ListenAddrs,
		ResourceManager,
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
		libp2p.ForceReachabilityPublic(),
	}
	return core.Option{
		LpOpt: opt,
		ID:    "",
	}
}

var ResourceManager = func(cfg *libp2p.Config) error {
	// Default memory limit: 1/8th of total memory, minimum 128MB, maximum 1GB
	limits := rcmgr.DefaultLimits
	libp2p.SetDefaultServiceLimits(&limits)
	limiter := rcmgr.NewFixedLimiter(rcmgr.InfiniteLimits)
	mgr, err := rcmgr.NewResourceManager(limiter)
	if err != nil {
		return err
	}

	return cfg.Apply(libp2p.ResourceManager(mgr))
}


func getNumFDs() int {
	var l unix.Rlimit
	if err := unix.Getrlimit(unix.RLIMIT_NOFILE, &l); err != nil {
		log.Errorw("failed to get fd limit", "error", err)
		return 0
	}
	return int(l.Cur)
}
