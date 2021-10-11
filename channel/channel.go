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

package channel

import (
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	pchannel "perun.network/go-perun/channel"
)

const (
	// OffIdentityLen is the length of an OffIdentity in byte.
	OffIdentityLen = 32
	// OnIdentityLen is the length of an OnIdentity in byte.
	OnIdentityLen = 32
	// NonceLen is the length of a Nonce in byte.
	NonceLen = 32
	// SigLen is the length of a Sig in byte.
	SigLen = 64
	// FIDLen is the length of an FID in byte.
	FIDLen = 32
)

// MaxBalance is the highest possible value for a Balance.
var MaxBalance = Balance(types.NewU128(*new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 128), big.NewInt(1))))

type (
	// Nonce makes a channels ID unique by providing randomness to the params.
	Nonce = [NonceLen]byte
	// ChannelId same as Peruns ChannelID.
	ChannelId = pchannel.ID
	// FundingId used to identify the funding of a participant in a channel.
	FundingId = [FIDLen]byte
	// Off-chain identity.
	OffIdentity = [OffIdentityLen]byte
	// On-chain identity.
	OnIdentity = [OnIdentityLen]byte
	// Version of a state as defined by go-perun.
	Version = uint64
	// ChallengeDuration as defined by go-perun.
	ChallengeDuration = uint64
	// Balance is the balance of an on- or off-chain Address.
	Balance = types.U128
	// Balances are multiple Balance.
	Balances = []Balance
	// Sig is an off-chain signature.
	Sig = [SigLen]byte
	// Sigs are multiple Sig.
	Sigs = []Sig

	// Params uniquely identify a channel.
	Params struct {
		// Nonce is the unique nonce of a channel.
		Nonce Nonce
		// Participants are the off-chain participants of the channel.
		Participants []OffIdentity
		// ChallengeDuration is the duration that disputes can be refuted in.
		ChallengeDuration ChallengeDuration
	}

	// State is the state of a channel.
	State struct {
		// Channel unique id of the channel.
		Channel ChannelId
		// Versoion of the state.
		Version Version
		// Balances of the participants.
		Balances Balances
		// Final whether or not this state is final.
		Final bool
	}

	// Withdrawal used by a participant to withdraw his on-chain funds.
	Withdrawal struct {
		// Channel is the channel from which to withdraw from.
		Channel ChannelId
		// Part is the participant woh wants to withdraw.
		Part OffIdentity
		// Receiver is the receiver of the withdrawal.
		Receiver OnIdentity
	}

	// Funding used to calculate FundingIds.
	Funding struct {
		// Channel the channel to fund.
		Channel ChannelId
		// Part the participant to fund.
		Part OffIdentity
	}

	// RegisteredState channel state that was registered on-chain.
	RegisteredState struct {
		// State of the channel.
		State State
		// Timeout is the duration that state can be refuted in.
		Timeout ChallengeDuration
		// Concluded whether the channel is concluded.
		Concluded bool
	}
)

// NewFunding returns a new Funding.
func NewFunding(id ChannelId, part OffIdentity) *Funding {
	return &Funding{id, part}
}

// Id calculates the funding id by encoding and hashing the receiver.
func (p Funding) Id() (FundingId, error) {
	var fid FundingId
	data, err := ScaleEncode(&p)
	if err != nil {
		return fid, errors.WithMessage(err, "calculating funding id")
	}
	return FundingId(crypto.Keccak256Hash(data)), nil
}
