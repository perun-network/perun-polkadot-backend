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
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/stretchr/testify/require"
)

const (
	ConfigFile = "chain.json"
	// PastBlocks defines how many blocks into the past the event subs should
	// query.
	PastBlocks = 100
)

type (
	ChainCfg struct {
		// ChainUrl is the url of the chain's RPC endpoint.
		ChainUrl string `json:"chain_url"`
		// Network is the network id.
		NetworkId substrate.NetworkId `json:"network_id"`
		// BlockTime is the block time of the node in seconds.
		BlockTimeSec uint32 `json:"block_time_sec"`
		BlockTime    time.Duration
	}

	Setup struct {
		ChainCfg

		Api *substrate.Api
	}
)

func NewSetup(t *testing.T, cfgDir string) *Setup {
	cfg := LoadChainCfg(t, cfgDir)
	api, err := substrate.NewApi(cfg.ChainUrl, cfg.NetworkId, PastBlocks)
	require.NoError(t, err)

	return &Setup{cfg, api}
}

func LoadChainCfg(t *testing.T, cfgDir string) ChainCfg {
	var cfg ChainCfg
	data, err := ioutil.ReadFile(path.Join(cfgDir, ConfigFile))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &cfg))
	cfg.BlockTime = time.Duration(cfg.BlockTimeSec) * time.Second
	return cfg
}
