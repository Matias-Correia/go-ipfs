package node

import (
	"context"
	"strconv"
	"time"

	"github.com/ipfs/go-bitswap"
	"github.com/ipfs/go-bitswap/network"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	config "github.com/ipfs/go-ipfs-config"
	exchange "github.com/ipfs/go-ipfs-exchange-interface"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"go.uber.org/fx"

	"github.com/ipfs/go-ipfs/core/node/helpers"

	logrpc "github.com/ipfs/go-bitswap/logrpc"
)

const (
	// Docs: https://github.com/ipfs/go-ipfs/blob/master/docs/config.md#internalbitswap
	DefaultEngineBlockstoreWorkerCount = 128
	DefaultTaskWorkerCount             = 8
	DefaultEngineTaskWorkerCount       = 8
	DefaultMaxOutstandingBytesPerPeer  = 1 << 20
	DefaultProviderMode				   = 1
	DefaultServerAddress			   = "localhost:50051"	
	//DefaultSessionAvgLatencyThreshold  = time.ParseDuration("300ms")
)

// OnlineExchange creates new LibP2P backed block exchange (BitSwap)
func OnlineExchange(cfg *config.Config, provide bool) interface{} {
	return func(mctx helpers.MetricsCtx, lc fx.Lifecycle, host host.Host, rt routing.Routing, bs blockstore.GCBlockstore) exchange.Interface {
		
		var internalBsCfg config.InternalBitswap
		if cfg.Internal.Bitswap != nil {
			internalBsCfg = *cfg.Internal.Bitswap
		}
		var serveraddr string = DefaultServerAddress
		if internalBsCfg.ServerAddress != "" {
			serveraddr = internalBsCfg.ServerAddress
		}
		
		gw := logrpc.New(serveraddr)

		bitswapNetwork := network.NewFromIpfsHost(host, rt, serveraddr, gw.GetChan())

		var providerSMode int = DefaultProviderMode
		if int(internalBsCfg.ProviderSelectionMode.WithDefault(DefaultProviderMode)) != 0{
			providerSMode = int(internalBsCfg.ProviderSelectionMode.WithDefault(DefaultProviderMode))
		}
		sessionavglatthreshold, _ := time.ParseDuration("300ms")
		if internalBsCfg.SessionAvgLatencyThreshold != 0{
			sessionavglatthreshold, _ = time.ParseDuration(strconv.Itoa(int(internalBsCfg.SessionAvgLatencyThreshold))+"ms")
		}

		opts := []bitswap.Option{
			bitswap.ProvideEnabled(provide),
			bitswap.EngineBlockstoreWorkerCount(int(internalBsCfg.EngineBlockstoreWorkerCount.WithDefault(DefaultEngineBlockstoreWorkerCount))),
			bitswap.TaskWorkerCount(int(internalBsCfg.TaskWorkerCount.WithDefault(DefaultTaskWorkerCount))),
			bitswap.EngineTaskWorkerCount(int(internalBsCfg.EngineTaskWorkerCount.WithDefault(DefaultEngineTaskWorkerCount))),
			bitswap.MaxOutstandingBytesPerPeer(int(internalBsCfg.MaxOutstandingBytesPerPeer.WithDefault(DefaultMaxOutstandingBytesPerPeer))),
			//bitswap.ProviderSelectionMode(int(internalBsCfg.ProviderSelectionMode.WithDefault(DefaultProviderMode)))
			//bitswap.ServerAddress(serveraddr)
		}
		exch := bitswap.New(helpers.LifecycleCtx(mctx, lc), bitswapNetwork, bs, providerSMode, serveraddr, sessionavglatthreshold, gw, opts...)
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return exch.Close()
			},
		})
		return exch

	}
}
