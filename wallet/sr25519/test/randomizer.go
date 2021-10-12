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
	ptest "perun.network/go-perun/wallet/test"

	pkgsr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
)

// Randomizer implements the wallet/test.Randomizer interface.
type Randomizer struct {
	wallet *Wallet
}

// NewRandomizer returns a new Randomizer.
func NewRandomizer() *Randomizer {
	return &Randomizer{NewWallet()}
}

// NewRandomAddress samples a random address from the passed entropy source.
func (*Randomizer) NewRandomAddress(rng *rand.Rand) pwallet.Address {
	pk, err := pkgsr25519.NewPkFromRng(rng)
	if err != nil {
		panic(err)
	}
	return sr25519.NewAddressFromPk(pk)
}

// RandomWallet returns the random wallet of the Randomizer.
func (r *Randomizer) RandomWallet() ptest.Wallet {
	return r.wallet
}

// NewWallet returns a new wallet.
func (*Randomizer) NewWallet() ptest.Wallet {
	return NewWallet()
}
