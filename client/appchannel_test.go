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
	"testing"
	"time"

	"perun.network/go-perun/channel"
	"perun.network/go-perun/client"
	clienttest "perun.network/go-perun/client/test"
	test "perun.network/go-perun/wallet/test"
	"perun.network/go-perun/wire"
	pkgtest "polycry.pt/poly-go/test"

	pchannel "github.com/perun-network/perun-polkadot-backend/channel"
	ptest "github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
)

func TestAppChannel(t *testing.T) {
	rng := pkgtest.Prng(t)
	s := ptest.NewSetup(t)
	setups := makeRoleSetups(rng, s, [2]string{"Paul", "Paula"})

	const A, B = 0, 1
	roles := [2]clienttest.Executer{
		clienttest.NewPaul(t, setups[A]),
		clienttest.NewPaula(t, setups[B]),
	}

	appAddress := test.NewRandomAddress(rng)
	app := channel.NewMockApp(appAddress)
	channel.RegisterApp(app)

	execConfig := &clienttest.ProgressionExecConfig{
		BaseExecConfig: clienttest.MakeBaseExecConfig(
			[2]wire.Address{setups[A].Identity.Address(), setups[B].Identity.Address()},
			pchannel.Asset,
			[2]*big.Int{big.NewInt(100000000000000), big.NewInt(100000000000000)},
			client.WithApp(app, channel.NewMockOp(channel.OpValid)),
		),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	clienttest.ExecuteTwoPartyTest(ctx, t, roles, execConfig)
}
