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
	"testing"

	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	wallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddress_SS58Addr tests that the SS58 address calculation returns the same
// result as polkadot.js.org.
func TestAddress_SS58Addr(t *testing.T) {
	for _, setup := range test.LoadDevAccounts("../../") {
		// Check AccountId
		gotId := wallet.AsAddr(setup.Acc.Address()).AccountId()
		assert.Equal(t, []byte(setup.Id), []byte(gotId[:]))

		// Check SS58 addresses
		for _, addr := range setup.Addr {
			gotAddr, err := substrate.SS58Address(wallet.AsAddr(setup.Acc.Address()).AccountId(), addr.Network)
			require.NoError(t, err)
			assert.Equal(t, addr.Value, gotAddr)
		}
	}
}
