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
	"time"

	"github.com/perun-network/perun-polkadot-backend/channel/pallet/test"
	sr25519test "github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"
	clienttest "perun.network/go-perun/client/test"
	"perun.network/go-perun/watcher/local"
	"perun.network/go-perun/wire"
)

func makeRoleSetups(s *test.Setup, names [2]string) (setup [2]clienttest.RoleSetup) {
	bus := wire.NewLocalBus()
	for i := 0; i < len(setup); i++ {
		watcher, err := local.NewWatcher(s.Adjs[i])
		if err != nil {
			panic(err)
		}
		setup[i] = clienttest.RoleSetup{
			Name:              names[i],
			Identity:          s.Accs[i].Acc,
			Bus:               bus,
			Funder:            s.Funders[i],
			Adjudicator:       s.Adjs[i],
			Wallet:            sr25519test.NewWallet(),
			Timeout:           TestTimeoutBlocks * time.Second,
			ChallengeDuration: 5, // 5 sec timeout
			Watcher:           watcher,
		}
	}
	return
}
