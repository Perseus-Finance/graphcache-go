package main

import (
	"fmt"

	"github.com/CrocSwap/graphcache-go/controller"
	"github.com/CrocSwap/graphcache-go/loader"
	"github.com/CrocSwap/graphcache-go/models"
	"github.com/CrocSwap/graphcache-go/server"
	"github.com/CrocSwap/graphcache-go/tables"
	"github.com/CrocSwap/graphcache-go/types"
	"github.com/CrocSwap/graphcache-go/views"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	netCfgPath := "../graphcache/webserver/config/networks.json"
	netCfg := loader.LoadNetworkConfig(netCfgPath)
	models := models.New()
	controller := controller.New(netCfg, models)

	goerlChainConfig, _ := netCfg["goerli"]
	goerliCntrl := controller.OnNetwork(types.NetworkName("goerli"))
	cfg := loader.SyncChannelConfig{
		Chain:   goerlChainConfig,
		Network: "goerli",
		Query:   "../graphcache/webserver/queries/balances.query",
	}

	tbl := tables.BalanceTable{}
	sync := loader.NewSyncChannel[tables.Balance, tables.BalanceSubGraph](
		tbl, cfg, goerliCntrl.IngestBalance)

	sync.SyncTableFromDb("../_data/database.db")
	sync.SyncTableToSubgraph()

	cfg.Query = "../graphcache/webserver/queries/swaps.query"
	tbl2 := tables.SwapsTable{}
	sync2 := loader.NewSyncChannel[tables.Swap, tables.SwapSubGraph](
		tbl2, cfg, func(l tables.Swap) { fmt.Println(l) })

	sync2.SyncTableFromDb("../_data/database.db")
	sync2.SyncTableToSubgraph()

	views := views.Views{Models: models}
	apiServer := server.APIWebServer{Views: &views}
	apiServer.Serve()
}
