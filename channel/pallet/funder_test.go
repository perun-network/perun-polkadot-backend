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

package pallet_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pchannel "perun.network/go-perun/channel"

	"github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
	chtest "github.com/perun-network/perun-polkadot-backend/channel/test"
)

func TestFunder_Fund(t *testing.T) {
	s := test.NewSetup(t)
	params, state := s.NewRandomParamAndState()
	fSetup := chtest.NewFundingSetup(params, state)

	err := test.FundAll(s.NewCtx(), s.Funders, fSetup.FReqs)
	require.NoError(t, err)

	// Check the on-chain balance.
	s.AssertDeposits(fSetup.Fids, fSetup.FinalBals)
}

// TestFunder_FundMultiple checks that funding twice results in twice the balance.
func TestFunder_FundMultiple(t *testing.T) {
	s := test.NewSetup(t)
	params, state := s.NewRandomParamAndState()
	fSetup := chtest.NewFundingSetup(params, state)

	err := test.FundAll(s.NewCtx(), s.Funders, fSetup.FReqs)
	require.NoError(t, err)
	// fund again
	err = test.FundAll(s.NewCtx(), s.Funders, fSetup.FReqs)
	require.NoError(t, err)

	// Check the on-chain balance.
	finalBals := test.Multiply(2, fSetup.FinalBals...)
	s.AssertDeposits(fSetup.Fids, finalBals)
}

func TestFunder_Timeout(t *testing.T) {
	s := test.NewSetup(t)
	params, state := s.NewRandomParamAndState()
	fSetup := chtest.NewFundingSetup(params, state)

	// Bob did not fund and times out.
	wantErr := makeTimeoutErr(1)
	// Only call Alice' funder.
	gotErr := s.Funders[0].Fund(s.NewCtx(), *fSetup.FReqs[0])
	// Check that the funder returned the correct error.
	assert.True(t, pchannel.IsFundingTimeoutError(gotErr))
	assert.Equal(t, wantErr.Error(), gotErr.Error()) // ErrorIs does not work here.
}

func makeTimeoutErr(idx pchannel.Index) error {
	return pchannel.NewFundingTimeoutError([]*pchannel.AssetFundingError{{Asset: 0, TimedOutPeers: []pchannel.Index{idx}}})
}
