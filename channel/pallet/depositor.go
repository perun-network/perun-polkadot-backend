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
	"context"

	"github.com/perun-network/perun-polkadot-backend/channel"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/pkg/errors"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/log"
	pwallet "perun.network/go-perun/wallet"
)

type (
	// Depositor deposits funds into a channel.
	Depositor struct {
		log.Embedding

		pallet *Pallet
	}

	// DepositReq contains arguments to specify a Deposit.
	DepositReq struct {
		Balance   pchannel.Bal
		Account   pwallet.Account
		FundingId channel.FundingId
	}
)

// NewDepositReq returns a new DepositReq.
func NewDepositReq(bal pchannel.Bal, acc pwallet.Account, fid channel.FundingId) *DepositReq {
	return &DepositReq{bal, acc, fid}
}

// WrapDepositReq wraps a perun funding request into a deposit request.
func WrapDepositReq(req *pchannel.FundingReq, acc pwallet.Account) (*DepositReq, error) {
	var fid channel.FundingId
	if !req.Agreement.Equal(req.State.Balances) && (len(req.Agreement) == 1) {
		return nil, errors.New("invalid funding request")
	}
	bal := req.Agreement[0][req.Idx]
	fid, err := channel.MakeFundingReqFromPerun(req).Id()
	return NewDepositReq(bal, acc, fid), err
}

// NewDepositor returns a new Depositor.
func NewDepositor(pallet *Pallet) *Depositor {
	return &Depositor{log.MakeEmbedding(log.Get()), pallet}
}

// Deposit deposits funds into the pallet for the funding id from the request.
func (d *Depositor) Deposit(ctx context.Context, req *DepositReq) error {
	ext, err := d.pallet.BuildDeposit(req.Account, req.Balance, req.FundingId)
	if err != nil {
		return err
	}
	d.Log().WithField("fid", req.FundingId).Debugf("Depositing %v", req.Balance)
	tx, err := d.pallet.Transact(ext)
	if err != nil {
		return err
	}
	defer tx.Close()
	return tx.WaitUntil(ctx, substrate.ExtIsFinal)
}
