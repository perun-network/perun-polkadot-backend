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

package test_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	pkgtest "polycry.pt/poly-go/test"

	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"
)

// TestRandomizer_RandomAddress tests that random addresses are distinct.
func TestRandomizer_RandomAddress(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewRandomizer()
	addr := r.NewRandomAddress(rng)

	for i := 0; i < 1000; i++ {
		addr2 := r.NewRandomAddress(rng)
		require.False(t, addr.Equals(addr2))
	}
}

// TestRandomWallet_RandomAccount tests that random accounts are distinct.
func TestRandomWallet_RandomAccount(t *testing.T) {
	rng := pkgtest.Prng(t)
	r := test.NewWallet()
	addr := r.NewRandomAccount(rng).Address()

	for i := 0; i < 1000; i++ {
		addr2 := r.NewRandomAccount(rng).Address()
		require.False(t, addr.Equals(addr2))
	}
}
