package main

import (
	"flag"

	"github.com/CrocSwap/graphcache-go/cache"
	"github.com/CrocSwap/graphcache-go/controller"
	"github.com/CrocSwap/graphcache-go/loader"
	"github.com/CrocSwap/graphcache-go/server"
	"github.com/CrocSwap/graphcache-go/views"
)

func main() {
	var netCfgPath = flag.String("netCfg", "./config/testnets.json", "network config file")
	var apiPath = flag.String("apiPath", "gcgo", "API server root path")
	var noRpcMode = flag.Bool("noRpcMode", false, "Run in mode with no RPC calls")
	flag.Parse()

	netCfg := loader.LoadNetworkConfig(*netCfgPath)
	onChain := loader.NewOnChainLoader(netCfg)

	cache := cache.New()
	cntrl := controller.New(netCfg, cache, onChain)

	if *noRpcMode {
		nonQuery := loader.NonCrocQuery{}
		cntrl = controller.NewOnQuery(netCfg, cache, &nonQuery)
	}

	for network, chainCfg := range netCfg {
		controller.NewSubgraphSyncer(cntrl, chainCfg, network)
	}

	views := views.Views{Cache: cache, OnChain: onChain}
	apiServer := server.APIWebServer{Views: &views}
	apiServer.Serve(*apiPath)
}
