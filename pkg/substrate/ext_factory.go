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
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

type (
	// ExtFactory can be used to build Extrinsics.
	ExtFactory struct {
		api *Api
	}

	// ExtName identifies an Extrinsic call by its name.
	ExtName struct {
		pallet, function string
	}
)

// NewExtName creates a new CallName for the given pallet and call name.
// Example: NewExtName("PerunModule", "deposit")
func NewExtName(pallet, function string) *ExtName {
	return &ExtName{pallet, function}
}

// NewExtFactory returns a new ExtrinsicBuilder.
func NewExtFactory(sub *Api) *ExtFactory {
	return &ExtFactory{sub}
}

// New returns a new Extrinsic with the given args.
func (b *ExtFactory) New(name *ExtName, args []interface{}) (*types.Extrinsic, error) {
	meta := b.api.Metadata()
	call, err := types.NewCall(meta, name.String(), args...)
	if err != nil {
		return nil, err
	}
	ext := types.NewExtrinsic(call)
	return &ext, nil
}

// SigOpts returns the default signature options for an address.
func (b *ExtFactory) SigOptions(addr types.AccountID) (*types.SignatureOptions, error) {
	info, err := b.api.AccountInfo(addr)
	if err != nil {
		return nil, err
	}
	genesis, err := b.api.BlockHash(0)
	if err != nil {
		return nil, err
	}
	runtime, err := b.api.RuntimeVersion()
	if err != nil {
		return nil, err
	}

	return &types.SignatureOptions{
		BlockHash:          genesis,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesis,
		Nonce:              types.NewUCompactFromUInt(uint64(info.Nonce)),
		SpecVersion:        runtime.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: runtime.TransactionVersion,
	}, nil
}

// String formats an ExtrinsicName in the form `Pallet.Function`.
func (f *ExtName) String() string {
	return fmt.Sprintf("%s.%s", f.pallet, f.function)
}
