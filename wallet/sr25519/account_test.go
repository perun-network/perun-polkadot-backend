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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"
)

// TestWallet_Verify tests that a constant signature that was generated at
// polkadot.js.org can be verified.
func TestWallet_Verify(t *testing.T) {
	for _, setup := range test.LoadDevAccounts(t, "../../") {
		backend := new(sr25519.Backend)
		ok, err := backend.VerifySignature(setup.Msg, setup.Sig, setup.Acc.Address())
		require.NoError(t, err)
		assert.True(t, ok)
	}
}
