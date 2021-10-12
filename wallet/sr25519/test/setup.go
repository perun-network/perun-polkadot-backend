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

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"perun.network/go-perun/wallet"

	pkgsr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
)

const (
	// file is the path to the config file which contains the dev accounts.
	file = "accounts.json"
)

type (
	// DevAccount contains an account for testing and developing which is
	// funded by the substrate node.
	DevAccount struct {
		Acc                wallet.Account
		Seed, Id, Msg, Sig Hex
		Addr               []SS58Addr
	}
	// SS58Addr contains an address and a network id.
	SS58Addr struct {
		Value   string
		Network substrate.NetworkId
	}
	Hex []byte
)

// UnmarshalJSON unmarshalls a Hex string from a byte slice.
func (s *Hex) UnmarshalJSON(data []byte) error {
	var hex string
	if err := json.Unmarshal(data, &hex); err != nil {
		return err
	}
	*s = hexutil.MustDecode(hex)
	return nil
}

// LoadDevAccounts loads all dev accounts by assuming that the config file is
// in the passed directory.
func LoadDevAccounts(t *testing.T, dir string) []*DevAccount {
	data, err := ioutil.ReadFile(path.Join(dir, file))
	require.NoError(t, err)

	var setups []*DevAccount
	err = json.Unmarshal(data, &setups)
	require.NoError(t, err)
	wallet := sr25519.NewWallet()
	for _, setup := range setups {
		sk, err := pkgsr25519.NewSk(setup.Seed)
		require.NoError(t, err)
		setup.Acc = wallet.ImportSK(sk)
	}
	return setups
}
