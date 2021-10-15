// Copyright 2021 PolyCrypt GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/stretchr/testify/require"
)

type (
	ChainCfg struct {
		// ChainUrl is the url of the chain's RPC endpoint.
		ChainUrl string `json:"chain_url"`
		// Network is the network ID.
		NetworkID substrate.NetworkId `json:"network_id"`
		// BlockTime is the block time of the node in seconds.
		BlockTimeSec uint32 `json:"block_time_sec"`
		BlockTime    time.Duration
	}

	Setup struct {
		ChainCfg

		Api *substrate.Api
	}
)

// PastBlocks defines how many blocks into the past the event subs should
// query.
const PastBlocks = 100

//go:embed chain.json
var chainCfgFile []byte

func NewSetup(t *testing.T) *Setup {
	cfg := LoadChainCfg(t)
	api, err := substrate.NewApi(cfg.ChainUrl, cfg.NetworkID, PastBlocks)
	require.NoError(t, err)

	return &Setup{cfg, api}
}

func LoadChainCfg(t *testing.T) ChainCfg {
	var cfg ChainCfg
	require.NoError(t, json.Unmarshal(chainCfgFile, &cfg))
	cfg.BlockTime = time.Duration(cfg.BlockTimeSec) * time.Second
	return cfg
}
