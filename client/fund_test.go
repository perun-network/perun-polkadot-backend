// Copyright 2022 - See NOTICE file for copyright holders.
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

package client_test

import (
	"context"
	"math/big"
	"math/rand"
	"testing"
	"time"

	dotchannel "github.com/perun-network/perun-polkadot-backend/channel"
	"github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
	"perun.network/go-perun/channel"
	ctest "perun.network/go-perun/client/test"
	pkgtest "polycry.pt/poly-go/test"
)

func TestFundRecovery(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	ctest.TestFundRecovery(
		ctx,
		t,
		ctest.FundSetup{
			ChallengeDuration: 1,
			FridaInitBal:      big.NewInt(100000000000000),
			FredInitBal:       big.NewInt(50000000000000),
			BalanceDelta:      big.NewInt(1000000000),
		},
		func(r *rand.Rand) ([2]ctest.RoleSetup, channel.Asset) {
			rng := pkgtest.Prng(t)
			s := test.NewSetup(t)
			roles := makeRoleSetups(rng, s, [2]string{"Frida", "Fred"})
			return roles, dotchannel.Asset
		},
	)
}
