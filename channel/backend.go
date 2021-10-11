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

package channel

import (
	"io"
	"log"

	eth "github.com/ethereum/go-ethereum/crypto"
	pchannel "perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"
)

// Backend implements the Backend interface.
type Backend struct{}

// NewBackend returns a new Backend.
func NewBackend() *Backend {
	return &Backend{}
}

// CalcID calculates the channelID.
func (*Backend) CalcID(p *pchannel.Params) (id pchannel.ID) {
	return CalcID(p)
}

// Sign signs a state with the passed account.
func (*Backend) Sign(acc pwallet.Account, state *pchannel.State) (pwallet.Sig, error) {
	data, err := ScaleEncode(NewState(state))
	if err != nil {
		return nil, err
	}
	return acc.SignData(data)
}

// Verify verifies a signature of a state.
func (*Backend) Verify(addr pwallet.Address, state *pchannel.State, sig pwallet.Sig) (bool, error) {
	data, err := ScaleEncode(NewState(state))
	if err != nil {
		return false, err
	}
	return pwallet.VerifySignature(data, sig, addr)
}

// DecodeAsset returns an instance of the Asset type. Cannot error.
func (*Backend) DecodeAsset(io.Reader) (pchannel.Asset, error) {
	return &Asset{}, nil
}

// CalcID calculates the channelID by encoding and hashing the params.
func CalcID(p *pchannel.Params) (id pchannel.ID) {
	bytes, err := ScaleEncode(NewParams(p))
	if err != nil {
		log.Panicf("could not encode parameters: %v", err)
	}
	return eth.Keccak256Hash(bytes)
}
