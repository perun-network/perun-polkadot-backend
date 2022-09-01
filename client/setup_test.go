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
	"math/rand"
	"time"

	"github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/perun-network/perun-polkadot-backend/wallet"
	sr25519test "github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"
	perunchannel "perun.network/go-perun/channel"
	clienttest "perun.network/go-perun/client/test"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
	netwire "perun.network/go-perun/wire/net/simple"
)

func makeRoleSetups(rng *rand.Rand, s *test.Setup, names [2]string) (setup [2]clienttest.RoleSetup) {
	bus := wire.NewLocalBus()
	for i := 0; i < len(setup); i++ {
		watcher, err := local.NewWatcher(s.Adjs[i])
		if err != nil {
			panic(err)
		}
		acc := wallet.AsAddr(s.Accs[i].Acc.Address())
		setup[i] = clienttest.RoleSetup{
			Name:              names[i],
			Identity:          netwire.NewRandomAccount(rng),
			Bus:               bus,
			Funder:            s.Funders[i],
			Adjudicator:       s.Adjs[i],
			Wallet:            sr25519test.NewWallet(),
			Timeout:           TestTimeoutBlocks * time.Second,
			ChallengeDuration: 5, // 5 sec timeout
			Watcher:           watcher,
			BalanceReader:     NewBalanceReader(s.API, acc),
		}
	}
	return
}

// BalanceReader is a balance reader used for testing. It is associated with a
// given account.
type BalanceReader struct {
	chain *substrate.API
	acc   wallet.Address
}

func NewBalanceReader(chain *substrate.API, acc wallet.Address) *BalanceReader {
	return &BalanceReader{
		chain: chain,
		acc:   acc,
	}
}

// Balance returns the asset balance of the associated account.
func (br *BalanceReader) Balance(_ perunchannel.Asset) perunchannel.Bal {
	accInfo, err := br.chain.AccountInfo(br.acc.AccountID())
	if err != nil {
		panic(err)
	}
	return accInfo.Free.Int
}
