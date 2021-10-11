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
	"github.com/centrifuge/go-substrate-rpc-client/v3/rpc/state"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/vedhavyas/go-subkey"
)

type (
	// NetworkId id of the substrate chain.
	NetworkId uint8

	// ExtrinsicSigner signs an Extrinsic by modifying it.
	ExtrinsicSigner interface {
		// SignExt signs the extrinsic with the specified options and network.
		SignExt(*gsrpc.Extrinsic, gsrpc.SignatureOptions, NetworkId) error
	}

	// StorageQueryer can be used to query the on-chain state.
	StorageQueryer interface {
		// QueryOne queries the given keys and return one result.
		// Errors if more than one or no result was found.
		QueryOne(pastBlocks gsrpc.BlockNumber, keys ...gsrpc.StorageKey) (*gsrpc.KeyValueOption, error)

		// Subscribe subscribes to the changes of a storage key.
		Subscribe(keys ...gsrpc.StorageKey) (*state.StorageSubscription, error)

		// StorageKey build a storage key.
		BuildKey(prefix, method string, args ...[]byte) (gsrpc.StorageKey, error)
	}

	// ChainReader is used to query on-chain state.
	ChainReader interface {
		// Metadata returns the latest metadata.
		Metadata() *gsrpc.Metadata
		// BlockHash returns the block hash for the given block number.
		BlockHash(gsrpc.BlockNumber) (gsrpc.Hash, error)
		// HeaderLatest returns the last header.
		HeaderLatest() *gsrpc.Header
	}
)

// SS58Address returns the SS58 of an Address for a specific network.
func SS58Address(addr gsrpc.AccountID, network NetworkId) (string, error) {
	return subkey.SS58Address(addr[:], uint8(network))
}
