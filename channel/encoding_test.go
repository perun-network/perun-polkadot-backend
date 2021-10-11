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

package channel_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pchtest "perun.network/go-perun/channel/test"
	pkgtest "perun.network/go-perun/pkg/test"

	// init backends + randomizer
	_ "github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"

	"github.com/perun-network/perun-polkadot-backend/channel"
	"github.com/perun-network/perun-polkadot-backend/channel/test"
)

func TestEncoding(t *testing.T) {
	rng := pkgtest.Prng(t)
	state := pchtest.NewRandomState(rng, test.DefaultRandomOpts())
	reg := &channel.RegisteredState{
		State:     *channel.NewState(state),
		Timeout:   123,
		Concluded: false,
	}
	// encode
	data, err := channel.ScaleEncode(reg)
	require.NoError(t, err)
	// decode
	dec := new(channel.RegisteredState)
	require.NoError(t, channel.ScaleDecode(dec, data))
	// equal?
	assert.Equal(t, reg, dec)
}
