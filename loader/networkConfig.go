package loader

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/CrocSwap/graphcache-go/types"
)

type ChainConfig struct {
	ChainID             int                 `json:"chain_id"`
	RPCs                map[string][]string `json:"rpcs"`
	Subgraph            string              `json:"subgraph"`
	DexContract         string              `json:"dex_contract"`
	QueryContract       string              `json:"query_contract"`
	QueryContractABI    string              `json:"query_contract_abi"`
	POAMiddleware       bool                `json:"poa_middleware"`
	BlockTime           float64             `json:"block_time"`
	Ignore              bool                `json:"ignore,omitempty"`
	EnableRPCCache      bool                `json:"enable_rpc_cache"`
	EnableSubgraphCache bool                `json:"enable_subgraph_cache"`
	KnockoutTickWidth   int                 `json:"knockout_ticks_width"`
}

type NetworkConfig map[types.NetworkName]ChainConfig

func LoadNetworkConfig(path string) NetworkConfig {
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var config NetworkConfig

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (c *NetworkConfig) ChainConfig(chainId types.ChainId) (ChainConfig, bool) {
	netName, isValid := c.NetworkForChainID(chainId)
	if isValid {
		cfg, hasCfg := (*c)[netName]
		if hasCfg {
			return cfg, true
		}
	}
	return ChainConfig{}, false
}

func (c *ChainConfig) RPCEndpoint() string {
	for _, rpcs := range c.RPCs {
		for _, rpc := range rpcs {
			return rpc
		}
	}
	log.Fatal("No configured RPC endpoint for " + types.IntToChainId(c.ChainID))
	return ""
}

func (c *NetworkConfig) NetworkForChainID(chainId types.ChainId) (types.NetworkName, bool) {
	for networkKey, configElem := range *c {
		if chainId == types.IntToChainId(configElem.ChainID) {
			return networkKey, true
		}
	}
	return "", false
}

func (c *NetworkConfig) ChainIDForNetwork(network types.NetworkName) (types.ChainId, bool) {
	lookup, ok := (*c)[network]
	if ok {
		return types.IntToChainId(lookup.ChainID), true
	} else {
		return "", false
	}
}
