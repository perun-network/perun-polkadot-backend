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
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/pkg/errors"
	pchannel "perun.network/go-perun/channel"
	pwallet "perun.network/go-perun/wallet"
)

// ScaleEncode encodes any struct according to the SCALE codec.
func ScaleEncode(obj interface{}) ([]byte, error) {
	return types.EncodeToBytes(obj)
}

// ScaleDecode decodes any struct according to the SCALE codec.
func ScaleDecode(obj interface{}, data []byte) error {
	return types.DecodeFromBytes(data, obj)
}

func MakeBalance(bal *big.Int) (Balance, error) {
	if bal.Sign() < 0 || bal.Cmp(MaxBalance.Int) > 0 {
		return Balance(types.NewU128(*big.NewInt(0))), errors.New("invalid balance")
	}
	return Balance(types.NewU128(*bal)), nil
}

func MakePerunBalance(bal Balance) *big.Int {
	return new(big.Int).Set(bal.Int)
}

func MakePerunAlloc(bals Balances) pchannel.Allocation {
	return pchannel.Allocation{
		Assets:   []pchannel.Asset{NewAsset()},
		Balances: MakePerunBalances(bals),
		Locked:   nil,
	}
}

func MakePerunBalances(bals Balances) pchannel.Balances {
	ret := make(pchannel.Balances, 1)
	ret[0] = make([]pchannel.Bal, len(bals))
	for i, bal := range bals {
		ret[0][i] = new(big.Int).Set(bal.Int)
	}
	return ret

}

func MakeChallengeDuration(challengeDuration uint64) ChallengeDuration {
	return ChallengeDuration(challengeDuration)
}

func MakePerunChallengeDuration(sec ChallengeDuration) uint64 {
	return uint64(sec)
}

func MakeDuration(sec ChallengeDuration) time.Duration {
	return time.Second * time.Duration(sec)
}

func MakeTime(sec ChallengeDuration) time.Time {
	return time.Unix(int64(sec), 0)
}

func MakeTimeout(sec ChallengeDuration, storage substrate.StorageQueryer) pchannel.Timeout {
	return substrate.NewTimeout(storage, MakeTime(sec))
}

func MakeNonce(nonce *big.Int) Nonce {
	if nonce.Sign() < 0 {
		panic("negative nonce")
	}
	if nonce.BitLen() > (8*NonceLen) || nonce.BitLen() == 0 {
		panic("empty nonce")
	}

	var ret Nonce
	copy(ret[:], nonce.Bytes())
	return ret
}

func NewState(s *pchannel.State) *State {
	if err := s.Valid(); err != nil {
		panic(err)
	}

	return &State{
		Channel:  s.ID,
		Version:  s.Version,
		Balances: MakeAlloc(&s.Allocation),
		Final:    s.IsFinal,
	}
}

func NewPerunState(s *State) *pchannel.State {
	return &pchannel.State{
		ID:         s.Channel,
		Version:    s.Version,
		Allocation: MakePerunAlloc(s.Balances),
		Data:       pchannel.NoData(),
		IsFinal:    s.Final,
	}
}

// MakeAlloc converts an Allocation to Balances since polkadot only
// supports one asset.
func MakeAlloc(a *pchannel.Allocation) Balances {
	if len(a.Assets) != 1 || len(a.Balances) != 1 {
		panic("expected exactly one asset")
	}
	if len(a.Locked) != 0 {
		panic("cannot lock balances")
	}
	var err error
	ret := make(Balances, len(a.Balances[0]))
	for i, bal := range a.Balances[0] {
		if ret[i], err = MakeBalance(bal); err != nil {
			panic(err)
		}
	}
	return ret
}

func NewWithdrawal(cid pchannel.ID, part, receiver pwallet.Address) *Withdrawal {
	return &Withdrawal{
		cid, MakeOffIdent(part), MakeOnIdent(receiver),
	}
}

func NewParams(p *pchannel.Params) *Params {
	ret := &Params{
		Nonce:             MakeNonce(p.Nonce),
		Participants:      MakeOffIdents(p.Parts),
		ChallengeDuration: MakeChallengeDuration(p.ChallengeDuration),
	}
	return ret
}

func MakeSig(sig pwallet.Sig) Sig {
	if len(sig) != SigLen {
		panic("wrong sig length")
	}
	var ret Sig
	copy(ret[:], sig)
	return ret
}

func MakeSigsFromPerun(sigs []pwallet.Sig) Sigs {
	ret := make([]Sig, len(sigs))
	for i, sig := range sigs {
		ret[i] = MakeSig(sig)
	}
	return ret
}

func MakeFundingReqFromPerun(req *pchannel.FundingReq) Funding {
	return Funding{
		req.State.ID,
		MakeOffIdent(req.Params.Parts[req.Idx]),
	}
}

func MakeOnIdent(addr pwallet.Address) OnIdentity {
	var ret OnIdentity
	data := addr.Bytes()
	if len(data) != OnIdentityLen {
		panic("wrong length")
	}
	copy(ret[:], data)
	return ret
}

func MakeOffIdent(part pwallet.Address) OffIdentity {
	var ret OffIdentity
	data := part.Bytes()
	if len(data) != OffIdentityLen {
		panic("wrong length")
	}
	copy(ret[:], data)
	return ret
}

func MakeOffIdents(parts []pwallet.Address) []OffIdentity {
	ret := make([]OffIdentity, len(parts))
	for i, part := range parts {
		ret[i] = MakeOffIdent(part)
	}
	return ret
}
