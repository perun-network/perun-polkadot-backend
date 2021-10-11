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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"perun.network/go-perun/wallet"

	pkgsr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
)

const (
	file = "accounts.json"
)

type (
	DevAccount struct {
		Acc                wallet.Account
		Seed, Id, Msg, Sig Hex
		Addr               []SS58Addr
	}
	SS58Addr struct {
		Value string
		// https://github.com/paritytech/substrate/blob/4c3a55e7ca5c4c85c1eb53fd82ed71029d952510/primitives/core/src/crypto.rs#L486
		Network substrate.NetworkId
	}
	Hex []byte
)

func (s *Hex) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return err
	}
	*s = hexutil.MustDecode(hex)
	return nil
}

func LoadDevAccounts(dir string) []*DevAccount {
	data, err := ioutil.ReadFile(path.Join(dir, file))
	must(err)

	var setups []*DevAccount
	must(json.Unmarshal(data, &setups))
	wallet := sr25519.NewWallet()
	for _, setup := range setups {
		sk, err := pkgsr25519.NewSk(setup.Seed)
		must(err)
		setup.Acc = wallet.ImportSK(sk)
	}
	return setups
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
