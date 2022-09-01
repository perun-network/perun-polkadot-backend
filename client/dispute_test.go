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

package client_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	pclient "perun.network/go-perun/client"
	clienttest "perun.network/go-perun/client/test"
	"perun.network/go-perun/wire"
	pkgtest "polycry.pt/poly-go/test"

	"github.com/perun-network/perun-polkadot-backend/channel"
	"github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
	"github.com/perun-network/perun-polkadot-backend/wallet"
)

const TestTimeoutBlocks = 100

func TestDisputeMalloryCarol(t *testing.T) {
	rng := pkgtest.Prng(t)
	s := test.NewSetup(t)
	ctx, cancel := context.WithTimeout(context.Background(), TestTimeoutBlocks*s.BlockTime)
	defer cancel()

	const (
		A, B = 0, 1 // Indices of Mallory and Carol
		// Maximal number of Extrinsics that a participant should send during the test.
		MaxNumExtSent = 6
	)
	var (
		name = [2]string{"Mallory", "Carol"}
		role [2]clienttest.Executer
	)

	setup := makeRoleSetups(rng, s, name)
	role[A] = clienttest.NewMallory(t, setup[A])
	role[B] = clienttest.NewCarol(t, setup[B])

	execConfig := &clienttest.MalloryCarolExecConfig{
		BaseExecConfig: clienttest.MakeBaseExecConfig(
			[2]wire.Address{setup[A].Identity.Address(), setup[B].Identity.Address()},
			channel.Asset,
			[2]*big.Int{big.NewInt(100000000000000), big.NewInt(100000000000000)},
			pclient.WithoutApp(),
		),
		NumPayments: [2]int{5, 0},
		TxAmounts:   [2]*big.Int{big.NewInt(2000000000000), big.NewInt(0)},
	}

	// Compensate for the fees of the extrinsics.
	epsilon := new(big.Int).SetUint64(test.DefaultExtFee * MaxNumExtSent)
	// Amount that will be send from Alice to Bob.
	aliceToBob := big.NewInt(int64(execConfig.NumPayments[A])*execConfig.TxAmounts[A].Int64() - int64(execConfig.NumPayments[B])*execConfig.TxAmounts[B].Int64())
	// Amount that will be send from Bob to Alice.
	bobToAlice := new(big.Int).Neg(aliceToBob)
	// Expected balance changes of the accounts.
	deltas := map[types.AccountID]*big.Int{
		wallet.AsAddr(s.Alice.Acc.Address()).AccountID(): aliceToBob,
		wallet.AsAddr(s.Bob.Acc.Address()).AccountID():   bobToAlice,
	}
	s.AssertBalanceChanges(deltas, epsilon, func() {
		clienttest.ExecuteTwoPartyTest(ctx, t, role, execConfig)
	})
}
