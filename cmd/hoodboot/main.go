package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	hoodboot "github.com/hood-chat/bootstrap-node"
	ds "github.com/ipfs/go-datastore"
	dsync "github.com/ipfs/go-datastore/sync"
	"github.com/ipfs/kubo/core/bootstrap"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	rcmgr "github.com/libp2p/go-libp2p/p2p/host/resource-manager"
	rh "github.com/libp2p/go-libp2p/p2p/host/routed"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr/net"
)

func main() {
	idPath := flag.String("id", "identity", "identity key file path")
	cfgPath := flag.String("config", "", "json configuration file; empty uses the default configuration")
	flag.Parse()

	cfg, err := hoodboot.LoadConfig(*cfgPath)
	if err != nil {
		panic(err)
	}
	privk, err := hoodboot.LoadIdentity(*idPath)
	if err != nil {
		panic(err)
	}

	var opts []libp2p.Option

	opts = append(opts,
		libp2p.UserAgent("hoodboot/1.0"),
		libp2p.Identity(privk),

		libp2p.EnableHolePunching(),
		libp2p.EnableNATService(),
		libp2p.ListenAddrStrings(cfg.Network.ListenAddrs...),
	)

	if len(cfg.Network.AnnounceAddrs) > 0 {
		var announce []ma.Multiaddr
		for _, s := range cfg.Network.AnnounceAddrs {
			a := ma.StringCast(s)
			announce = append(announce, a)
		}
		opts = append(opts,
			libp2p.AddrsFactory(func([]ma.Multiaddr) []ma.Multiaddr {
				return announce
			}),
		)
	} else {
		opts = append(opts,
			libp2p.AddrsFactory(func(addrs []ma.Multiaddr) []ma.Multiaddr {
				announce := make([]ma.Multiaddr, 0, len(addrs))
				for _, a := range addrs {
					if manet.IsPublicAddr(a) {
						announce = append(announce, a)
					}
				}
				return announce
			}),
		)
	}

	if cfg.Relay.Enable {
		opts = append(opts,
			libp2p.EnableRelay(),
		)
	} else {
		opts = append(opts,
			libp2p.DisableRelay(),
		)
	}

	if cfg.HolePunch.Enable {
		opts = append(opts,
			libp2p.EnableHolePunching(),
		)
	}

	cm, err := connmgr.NewConnManager(
		cfg.ConnMgr.ConnMgrLo,
		cfg.ConnMgr.ConnMgrHi,
		connmgr.WithGracePeriod(cfg.ConnMgr.ConnMgrGrace),
	)
	if err != nil {
		panic(err)
	}

	opts = append(opts,
		libp2p.ConnectionManager(cm),
	)

	rm, err := rcmgr.NewResourceManager(rcmgr.NewFixedLimiter(rcmgr.InfiniteLimits))
	if err != nil {
		panic(err)
	}

	opts = append(opts,
		libp2p.ResourceManager(rm),
	)

	host, err := libp2p.New(opts...)
	if err != nil {
		panic(err)
	}

	fmt.Printf("I am %s\n", host.ID())
	fmt.Printf("Public Addresses:\n")
	for _, addr := range host.Addrs() {
		fmt.Printf("\t%s/p2p/%s\n", addr, host.ID())
	}

	var kDht *dht.IpfsDHT = nil
	if cfg.DHT.Enable {
		// Setup DHT
		dstore := dsync.MutexWrap(ds.NewMapDatastore())

		// Make the DHT
		kDht = dht.NewDHT(context.Background(), host, dstore)

		// Make the routed host
		rh.Wrap(host, kDht)
	}

	bts, err := hoodboot.ParsePeers(cfg.Bootstrap.Peers)
	if err != nil {
		panic(err)
	}
	btconf := bootstrap.BootstrapConfigWithPeers(bts)
	btconf.MinPeerThreshold = cfg.Bootstrap.MinPeerThreshold

	// connect to the chosen ipfs nodes
	_, err = bootstrap.Bootstrap(host.ID(), host, kDht, btconf)
	if err != nil {
		panic(err)
	}

	go listenPprof(cfg.Daemon.PprofPort)
	time.Sleep(10 * time.Millisecond)

	select {}
}

func listenPprof(p int) {
	if p == -1 {
		fmt.Printf("The pprof debug is disabled\n")
		return
	}
	addr := fmt.Sprintf("localhost:%d", p)
	fmt.Printf("Registering pprof debug http handler at: http://%s/debug/pprof/\n", addr)
	switch err := http.ListenAndServe(addr, nil); err {
	case nil:
		// all good, server is running and exited normally.
	case http.ErrServerClosed:
		// all good, server was shut down.
	default:
		// error, try another port
		fmt.Printf("error registering pprof debug http handler at: %s: %s\n", addr, err)
		panic(err)
	}
}
