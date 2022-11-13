package main

import (
	"fmt"

	"github.com/hood-chat/core"
	"github.com/hood-chat/core/entity"
	"github.com/hood-chat/core/repo"
	"github.com/hood-chat/core/store"
	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/config"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/multiformats/go-multiaddr"
)

func init() {
	fmt.Println("Hello! init() function")
}

// Main function
func main() {
	err := logging.SetLogLevel("*", "DEBUG")
	if err != nil {
		panic(err)
	}
	s, err := store.NewStore("./data")
	if err != nil {
		panic(err)
	}
	rIdentity := repo.NewIdentityRepo(s)
	id, err := rIdentity.Get()
	if err != nil {
		id, err = entity.CreateIdentity("bootstraper")
		if err != nil {
			panic(err)
		}
	}
	opt := Option()

	opt.SetIdentity(&id)
	hb := core.DefaultRoutedHost{}
	if err != nil {
		panic(err)
	}

	_, err = hb.Create(opt)
	if err != nil {
		panic(err)
	}

	fmt.Println("Welcome to main() function")

	select {} // block forever
}

var ListenAddrs = func(cfg *config.Config) error {
	defaultIP4ListenAddr, err := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4001")
	if err != nil {
		return err
	}

	defaultIP6ListenAddr, err := multiaddr.NewMultiaddr("/ip6/::/tcp/4001")
	if err != nil {
		return err
	}
	return cfg.Apply(libp2p.ListenAddrs(
		defaultIP4ListenAddr,
		defaultIP6ListenAddr,
	))
}

func Option() core.Option {
	// Now, normally you do not just want a simple host, you want
	// that is fully configured to best support your p2p application.
	// Let's create a second host setting some more options.
	// Set your own keypair

	con, err := connmgr.NewConnManager(10, 200)
	if err != nil {
		panic(err)
	}

	opt := []libp2p.Option{
		libp2p.DefaultTransports,
		libp2p.DefaultSecurity,
		// Use the keypair we generated
		// Multiple listen addresses
		ListenAddrs,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(con),
		// libp2p.DefaultMuxers,
		// Let this host use relays and advertise itself on relays if
		// it finds it is behind NAT. Use libp2p.Relay(options...) to
		// enable active relays and more.
		// libp2p.EnableAutoRelay(),
		libp2p.EnableAutoRelay(),
		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
	}
	return core.Option{
		LpOpt: opt,
		ID:    "",
	}
}
