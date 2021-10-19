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

package sr25519_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	pkgtest "perun.network/go-perun/pkg/test"
	ptest "perun.network/go-perun/wallet/test"

	pkgsr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	wallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
)

func TestWallet_GenericSignatureSizeTest(t *testing.T) {
	s := newSetup(t, pkgtest.Prng(t))
	ptest.GenericSignatureSizeTest(t, s)
}

func newSetup(t require.TestingT, rng *rand.Rand) *ptest.Setup {
	w := wallet.NewWallet()

	sk1, err := pkgsr25519.NewSKFromRng(rng)
	require.NoError(t, err)
	acc1 := w.ImportSK(sk1)

	sk2, err := pkgsr25519.NewSKFromRng(rng)
	require.NoError(t, err)
	acc2 := w.ImportSK(sk2)

	data := make([]byte, rng.Intn(1024))
	_, err = rng.Read(data)
	require.NoError(t, err)

	zeroPK, _ := pkgsr25519.NewPK([]byte{})

	return &ptest.Setup{
		AddressInWallet: acc1.Address(),
		ZeroAddress:     wallet.NewAddressFromPK(zeroPK),
		Backend:         new(wallet.Backend),
		Wallet:          w,
		AddressEncoded:  acc2.Address().Bytes(),
		DataToSign:      data,
	}
}
