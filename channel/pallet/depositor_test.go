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

	"github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
	chtest "github.com/perun-network/perun-polkadot-backend/channel/test"
	"github.com/stretchr/testify/require"
)

func TestDepositor_Deposit(t *testing.T) {
	s := test.NewSetup(t)
	params, state := s.NewRandomParamAndState()
	fSetup := chtest.NewFundingSetup(params, state)
	dSetup := chtest.NewDepositSetup(fSetup, s.Alice.Acc, s.Bob.Acc)

	err := test.DepositAll(s.NewCtx(), s.Deps, dSetup.DReqs)
	require.NoError(t, err)

	// Check the on-chain balance.
	s.AssertDeposits(fSetup.Fids, fSetup.FinalBals)
}

func TestDepositor_DepositMultiple(t *testing.T) {
	s := test.NewSetup(t)
	params, state := s.NewRandomParamAndState()
	fSetup := chtest.NewFundingSetup(params, state)
	dSetup := chtest.NewDepositSetup(fSetup, s.Alice.Acc, s.Bob.Acc)

	err := test.DepositAll(s.NewCtx(), s.Deps, dSetup.DReqs)
	require.NoError(t, err)
	err = test.DepositAll(s.NewCtx(), s.Deps, dSetup.DReqs)
	require.NoError(t, err)

	// Check the on-chain balance.
	finalBals := test.Multiply(2, fSetup.FinalBals...)
	s.AssertDeposits(fSetup.Fids, finalBals)
}
