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

package substrate

import (
	"math/big"
	"sync"

	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v3"
	"github.com/centrifuge/go-substrate-rpc-client/v3/rpc/chain"
	"github.com/centrifuge/go-substrate-rpc-client/v3/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/pkg/errors"
	"perun.network/go-perun/log"
)

type (
	// Wraps a gsrpc.SubstrateAPI in a thread-safe way.
	// Exposes basic query functionality.
	Api struct {
		log.Embedding
		mtx     *sync.Mutex // protects all
		url     string
		api     *gsrpc.SubstrateAPI
		network NetworkId
		meta    *types.Metadata
	}
)

var (
	// ErrWrongNodeVersion returned if an invalid substrate node version was
	// detected.
	ErrWrongNodeVersion = errors.New("wrong node version")
)

// NewApi creates a new `Api` object. Can be retried in the error case.
func NewApi(url string, network NetworkId, pastBlocks types.BlockNumber) (*Api, error) {
	logger := log.MakeEmbedding(log.Get())
	log.WithField("url", url).Debugf("Connecting to node")
	api, err := gsrpc.NewSubstrateAPI(url)
	if err != nil {
		return nil, errors.WithMessage(err, "connecting to substrate node")
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if !meta.IsMetadataV13 {
		return nil, ErrWrongNodeVersion
	}
	return &Api{logger, new(sync.Mutex), url, api, network, meta}, err
}

// AccountInfo returns the account info for an Address. Can be used to get
// the free balance, nonce, and other.
func (a *Api) AccountInfo(addr types.AccountID) (AccountInfo, error) {
	key, err := types.CreateStorageKey(a.Metadata(), "System", "Account", addr[:])
	if err != nil {
		return AccountInfo{}, err
	}

	var info AccountInfo
	ok, err := a.StorageLatest(key, &info)
	if !ok {
		return AccountInfo{}, errors.Errorf("account not found: 0x%x", addr)
	}
	return info, err
}

// Metadata returns the metadata of the chain. The value is cached on startup.
func (a *Api) Metadata() *types.Metadata {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.meta
}

// Network returns the id of the network that the api is connected to.
// The value is cached on startup.
func (a *Api) Network() NetworkId {
	return a.network
}

// BlockHash returns the hash for the given block number.
func (a *Api) BlockHash(n uint64) (types.Hash, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.Chain.GetBlockHash(n)
}

// PastBlock queries `pastBlocks` into the past and returns the hash of the
// block.
func (a *Api) PastBlock(pastBlocks types.BlockNumber) (types.Hash, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	head, err := a.api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return types.Hash{}, err
	}
	blockNum := types.BlockNumber(0)
	if head.Number > pastBlocks {
		blockNum = head.Number - pastBlocks
	}
	return a.api.RPC.Chain.GetBlockHash(uint64(blockNum))
}

// RuntimeVersion queries and returns the last runtime version.
func (a *Api) RuntimeVersion() (*types.RuntimeVersion, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.State.GetRuntimeVersionLatest()
}

// Subscribe subscribes to multiple storage keys.
func (a *Api) Subscribe(keys ...types.StorageKey) (*state.StorageSubscription, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.State.SubscribeStorageRaw(keys)
}

// BuildKey builds a storage key.
func (a *Api) BuildKey(pallet, variable string, args ...[]byte) (types.StorageKey, error) {
	return types.CreateStorageKey(a.Metadata(), pallet, variable, args...)
}

// QueryAll returns all entries for `keys` from `startBlock` to the last block.
func (a *Api) QueryAll(keys []types.StorageKey, startBlock types.Hash) ([]types.StorageChangeSet, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.State.QueryStorageLatest(keys, startBlock)
}

// QueryOne queries the storage and expects to read at least one value.
// Returns the latest value or an error if none was found.
func (a *Api) QueryOne(pastBlocks types.BlockNumber, keys ...types.StorageKey) (*types.KeyValueOption, error) {
	firstBlock, err := a.PastBlock(pastBlocks)
	if err != nil {
		return nil, err
	}
	sets, err := a.QueryAll(keys, firstBlock)
	if err != nil {
		return nil, err
	}
	if len(sets) == 0 {
		return nil, errors.New("nothing found")
	}
	set := sets[len(sets)-1]
	if len(set.Changes) != 1 {
		return nil, errors.New("not exactly one change")
	}
	return &set.Changes[len(set.Changes)-1], nil
}

// LastHeader returns the last header.
func (a *Api) LastHeader() (*types.Header, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.Chain.GetHeaderLatest()
}

// SubscribeHeaders subscribes to new headers.
func (a *Api) SubscribeHeaders() (*chain.NewHeadsSubscription, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.Chain.SubscribeNewHeads()
}

// StorageLatest queries the latest storage state.
func (a *Api) StorageLatest(key types.StorageKey, target interface{}) (bool, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()

	return a.api.RPC.State.GetStorageLatest(key, target)
}

// SubmitAndWatch sends an Extrinsic and returns an Sub to listen for updates
// of the Extrinsic.
func (a *Api) SubmitAndWatch(ext *types.Extrinsic) (*ExtStatusSub, error) {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.Log().Debugf("Sending tx with nonce = %v", ((*big.Int)(&ext.Signature.Nonce)).Int64())
	sub, err := a.api.RPC.Author.SubmitAndWatchExtrinsic(*ext)
	return NewExtStatusSub(sub), err
}
