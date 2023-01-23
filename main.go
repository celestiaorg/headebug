package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/celestiaorg/celestia-node/header"
	"github.com/celestiaorg/celestia-node/libs/header/sync"
	"github.com/celestiaorg/celestia-node/logs"
	"github.com/celestiaorg/celestia-node/nodebuilder"
	modfraud "github.com/celestiaorg/celestia-node/nodebuilder/fraud"
	modheader "github.com/celestiaorg/celestia-node/nodebuilder/header"
	"github.com/celestiaorg/celestia-node/nodebuilder/node"
	"github.com/celestiaorg/celestia-node/nodebuilder/p2p"
	"github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"
	ma "github.com/multiformats/go-multiaddr"
	"go.uber.org/fx"
)

func main() {
	ctx := context.Background()
	err := run(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	timeout := time.Hour * 1

	logs.SetAllLoggers(log.LevelInfo)
	log.SetLogLevel("header/p2p", "debug")
	log.SetLogLevel("basichost", "debug")
	log.SetLogLevel("swarm2", "debug")

	bootstrappers := p2p.Bootstrappers{
		mustDecode("/ip4/207.154.220.138/udp/2125/quic/p2p/12D3KooWN7vfV1zZfbbf1ty1wszTkveEXdTu6F5cSBUkePjGcARD"),
	}
	tp := node.Light
	hcfg := modheader.DefaultConfig(tp)
	pcfg := p2p.DefaultConfig()
	store := nodebuilder.NewMemStore()
	app := fx.New(
		fx.StartTimeout(timeout),
		fx.Provide(context.Background),
		fx.Supply(tp),
		fx.Supply(p2p.Arabica),
		fx.Provide(store.Datastore),
		fx.Provide(store.Keystore),
		fx.Supply(bootstrappers),
		p2p.ConstructModule(tp, &pcfg),
		modheader.ConstructModule(tp, &hcfg),
		modfraud.ConstructModule(tp),
		fx.Invoke(func(sync *sync.Syncer[*header.ExtendedHeader]) {
			// invoke
		}),
	)

	err := app.Err()
	if err != nil {
		return err
	}

	app.Run()
	return nil
}

func mustDecode(s string) peer.AddrInfo {
	id, err := peer.AddrInfoFromP2pAddr(ma.StringCast(s))
	if err != nil {
		panic(err)
	}
	return *id
}
