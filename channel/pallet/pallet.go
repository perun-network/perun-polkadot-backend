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

package pallet

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/log"
	pwallet "perun.network/go-perun/wallet"

	"github.com/perun-network/perun-polkadot-backend/channel"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/perun-network/perun-polkadot-backend/wallet"
)

const (
	PALLET_PERUN = "PerunModule"
)

var (
	DEPOSIT          = substrate.NewExtName(PALLET_PERUN, "deposit")
	DISPUTE          = substrate.NewExtName(PALLET_PERUN, "dispute")
	CONCLUDE         = substrate.NewExtName(PALLET_PERUN, "conclude")
	CONCLUDE_DISPUTE = substrate.NewExtName(PALLET_PERUN, "conclude_dispute")
	WITHDRAW         = substrate.NewExtName(PALLET_PERUN, "withdraw")
)

// Pallet exposes all functions of the perun-pallet.
// https://github.com/perun-network/perun-polkadot-pallet
type Pallet struct {
	log.Embedding
	*substrate.Pallet
	meta *types.Metadata
}

// NewPerunPallet is a wrapper around `NewPallet` and calls it with the values
// for the perun-pallet.
func NewPerunPallet(api *substrate.Api) *substrate.Pallet {
	return substrate.NewPallet(api, PALLET_PERUN)
}

// NewPallet returns a new generic Pallet.
func NewPallet(pallet *substrate.Pallet, meta *types.Metadata) *Pallet {
	return &Pallet{log.MakeEmbedding(log.Get()), pallet, meta}
}

// Subscribe returns an EventSub that listens on all events of the pallet.
func (p *Pallet) Subscribe(f EventPredicate, pastBlocks types.BlockNumber) (*EventSub, error) {
	source, err := p.Pallet.Subscribe(pastBlocks)
	if err != nil {
		return nil, err
	}
	return NewEventSub(source, p.meta, f), nil
}

// QueryDeposit returns the deposit for a funding id or
// nil if no deposit was found.
func (p *Pallet) QueryDeposit(fid channel.FundingId, storage substrate.StorageQueryer, pastBlocks types.BlockNumber) (*channel.Balance, error) {
	log := p.Log().WithField("fid", fid)
	key, err := p.BuildQuery("Deposits", fid[:])
	if err != nil {
		return nil, err
	}
	log.Trace("Querying deposit")
	res, err := storage.QueryOne(pastBlocks, key)
	if err != nil {
		return nil, err
	}
	if len(res.StorageData) == 0 {
		return nil, nil
	}
	ret := new(channel.Balance)
	return ret, channel.ScaleDecode(ret, res.StorageData)
}

// QueryStateRegister returns the registered state of a channel
// or nil if no state was found.
func (p *Pallet) QueryStateRegister(cid channel.ChannelId, storage substrate.StorageQueryer, pastBlocks types.BlockNumber) (*channel.RegisteredState, error) {
	log := p.Log().WithField("cid", cid)
	key, err := p.BuildQuery("StateRegister", cid[:])
	if err != nil {
		return nil, err
	}
	log.Trace("Querying stateRegister")
	res, err := storage.QueryOne(pastBlocks, key)
	if err != nil {
		return nil, err
	}
	if len(res.StorageData) == 0 {
		return nil, nil
	}
	ret := new(channel.RegisteredState)
	return ret, channel.ScaleDecode(ret, res.StorageData)
}

// BuildDeposit returns an extrinsic that funds the specified funding id.
func (p *Pallet) BuildDeposit(acc pwallet.Account, _amount pchannel.Bal, fid channel.FundingId) (*types.Extrinsic, error) {
	amount, err := channel.MakeBalance(_amount)
	if err != nil {
		return nil, err
	}
	return p.BuildExt(DEPOSIT, []interface{}{
		fid,
		amount},
		wallet.AsAddr(acc.Address()).AccountId(),
		wallet.AsAcc(acc))
}

// BuildDispute returns an extrinsic that disputes a channel.
func (p *Pallet) BuildDispute(acc pwallet.Account, _params *pchannel.Params, _state *pchannel.State, _sigs []pwallet.Sig) (*types.Extrinsic, error) {
	return p.BuildExt(DISPUTE,
		[]interface{}{
			channel.NewParams(_params),
			channel.NewState(_state),
			channel.MakeSigsFromPerun(_sigs)},
		wallet.AsAddr(acc.Address()).AccountId(),
		wallet.AsAcc(acc))
}

// BuildConclude returns an extrinsic that concludes a channel.
func (p *Pallet) BuildConclude(acc pwallet.Account, _params *pchannel.Params, _state *pchannel.State, _sigs []pwallet.Sig) (*types.Extrinsic, error) {
	return p.BuildExt(CONCLUDE,
		[]interface{}{
			channel.NewParams(_params),
			channel.NewState(_state),
			channel.MakeSigsFromPerun(_sigs)},
		wallet.AsAddr(acc.Address()).AccountId(),
		wallet.AsAcc(acc))
}

// BuildWithdraw returns an extrinsic that withdraws all funds from the channel.
func (p *Pallet) BuildWithdraw(onChain, offChain pwallet.Account, cid pchannel.ID) (*types.Extrinsic, error) {
	part := offChain.Address()
	receiver := onChain.Address()

	withdrawal := channel.NewWithdrawal(cid, part, receiver)
	data, err := channel.ScaleEncode(withdrawal)
	if err != nil {
		return nil, err
	}
	sig, err := offChain.SignData(data)
	if err != nil {
		return nil, err
	}

	return p.BuildExt(WITHDRAW,
		[]interface{}{
			withdrawal,
			channel.MakeSig(sig)},
		wallet.AsAddr(onChain.Address()).AccountId(),
		wallet.AsAcc(onChain))
}
