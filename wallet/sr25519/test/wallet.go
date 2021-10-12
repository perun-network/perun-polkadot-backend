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
	"math/rand"

	pwallet "perun.network/go-perun/wallet"

	"github.com/ChainSafe/go-schnorrkel"
	pkgsr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
)

// Wallet is used for testing and implements the Wallet interface.
type Wallet struct {
	*sr25519.Wallet
}

// NewWallet returns a new Wallet.
func NewWallet() *Wallet {
	return &Wallet{sr25519.NewWallet()}
}

// NewRandomAccount samples a random account from the passed entropy source.
// Imports the account into the wallet and returns it.
func (w *Wallet) NewRandomAccount(rng *rand.Rand) pwallet.Account {
	sk, err := pkgsr25519.NewSkFromRng(rng)
	if err != nil {
		panic(err)
	}
	return w.Wallet.ImportSK(sk)
}

// NewAddressZero returns a zero address that is strictly smaller than all
// other addresses.
func NewAddressZero() *sr25519.Address {
	return sr25519.NewAddressFromPk(ZeroPk())
}

// ZeroPk returns a PK that can be used to create a zero address.
func ZeroPk() *schnorrkel.PublicKey {
	zero, _ := pkgsr25519.NewPk([]byte{})
	return zero
}
