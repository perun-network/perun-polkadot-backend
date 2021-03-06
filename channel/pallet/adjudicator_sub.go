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
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/perun-network/perun-polkadot-backend/channel"
	pkg_sr25519 "github.com/perun-network/perun-polkadot-backend/pkg/sr25519"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	pchannel "perun.network/go-perun/channel"
	"perun.network/go-perun/log"
	pkgsync "polycry.pt/poly-go/sync"
)

// AdjudicatorSub implements the AdjudicatorSubscription interface.
type AdjudicatorSub struct {
	*pkgsync.Closer
	log.Embedding

	cid     channel.ChannelID
	sub     *EventSub
	pallet  *Pallet
	storage substrate.StorageQueryer
	err     chan error
}

// NewAdjudicatorSub returns a new AdjudicatorSub. Will return all events from
// the `pastBlocks` past blocks and all events from future blocks.
func NewAdjudicatorSub(cid channel.ChannelID, p *Pallet, storage substrate.StorageQueryer, pastBlocks types.BlockNumber) (*AdjudicatorSub, error) {
	sub, err := p.Subscribe(isAdjEvent(cid), pastBlocks)
	if err != nil {
		return nil, err
	}
	ret := &AdjudicatorSub{new(pkgsync.Closer), log.MakeEmbedding(log.Default()), cid, sub, p, storage, make(chan error, 1)}
	ret.OnCloseAlways(func() {
		if err := ret.sub.Close(); err != nil {
			ret.Log().WithError(err).Error("Could not close Closer.")
		}
		close(ret.err)
	})
	return ret, nil
}

// isAdjEvent returns a predicate that decides whether or not an event is
// relevant for the adjudicator and concerns a specific channel.
func isAdjEvent(cid channel.ChannelID) EventPredicate {
	return func(e channel.PerunEvent) bool {
		return channel.EventIsDisputed(cid)(e) || channel.EventIsProgressed(cid)(e) || channel.EventIsConcluded(cid)(e)
	}
}

// Next implements the AdjudicatorSub.Next function.
func (s *AdjudicatorSub) Next() pchannel.AdjudicatorEvent {
	if s.IsClosed() {
		return nil
	}
	var last channel.PerunEvent
	// Wait for event or closed.
	select {
	case event := <-s.sub.Events():
		if channel.EventIsDisputed(s.cid)(event) || channel.EventIsProgressed(s.cid)(event) || channel.EventIsConcluded(s.cid)(event) {
			last = event
		}
	case <-s.Closed():
		return nil
	}

loop:
	for {
		select {
		case event := <-s.sub.Events():
			if channel.EventIsDisputed(s.cid)(event) || channel.EventIsProgressed(s.cid)(event) || channel.EventIsConcluded(s.cid)(event) {
				last = event
			}
		case <-s.Closed():
			return nil
		// wait until there is no new event for a specific time
		case <-time.After(100 * time.Millisecond):
			break loop
		}
	}
	// Convert event.
	event, err := s.makePerunEvent(last)
	if err != nil {
		s.err <- err
		if err := s.Closer.Close(); err != nil {
			s.Log().WithError(err).Error("Could not close Closer.")
		}
		return nil
	}
	return event
}

// makePerunEvent creates a Perun event from a generic event.
func (s *AdjudicatorSub) makePerunEvent(event channel.PerunEvent) (pchannel.AdjudicatorEvent, error) {
	switch event := event.(type) {
	case *channel.DisputedEvent:
		s.Log().Trace("AdjudicatorSub creating DisputedEvent")
		dispute, err := s.pallet.QueryStateRegister(event.Cid, s.storage, 0)
		if err != nil {
			return nil, err
		}

		return &pchannel.RegisteredEvent{
			AdjudicatorEventBase: pchannel.AdjudicatorEventBase{
				IDV:      event.Cid,
				VersionV: event.State.Version,
				TimeoutV: channel.MakeTimeout(dispute.Timeout, s.storage),
			},
			State: nil, // only needed for virtual channel support
			Sigs:  nil, // only needed for virtual channel support
		}, nil
	case *channel.ProgressedEvent:
		s.Log().Trace("AdjudicatorSub creating ProgressedEvent")
		dispute, err := s.pallet.QueryStateRegister(event.Cid, s.storage, 0)
		if err != nil {
			return nil, err
		}

		appPK, err := pkg_sr25519.NewPK(event.App[:])
		if err != nil {
			return nil, err
		}
		appAddr := sr25519.NewAddressFromPK(appPK)
		app, err := pchannel.Resolve(appAddr)
		if err != nil {
			return nil, err
		}

		return &pchannel.ProgressedEvent{
			AdjudicatorEventBase: pchannel.AdjudicatorEventBase{
				IDV:      event.Cid,
				VersionV: event.Version,
				TimeoutV: channel.MakeTimeout(dispute.Timeout, s.storage),
			},
			State: channel.NewPerunState(&dispute.State, app),
		}, nil
	case *channel.ConcludedEvent:
		s.Log().Trace("AdjudicatorSub creating ConcludedEvent")
		dispute, err := s.pallet.QueryStateRegister(event.Cid, s.storage, 0)
		if err != nil {
			return nil, err
		}

		return &pchannel.ConcludedEvent{
			AdjudicatorEventBase: pchannel.AdjudicatorEventBase{
				IDV:      event.Cid,
				VersionV: dispute.State.Version,
				TimeoutV: channel.MakeTimeout(dispute.Timeout, s.storage),
			},
		}, nil
	default:
		panic(fmt.Sprintf("unknown event: %#v", event))
	}
}

// Err implements the AdjudicatorSub.Err function.
func (s *AdjudicatorSub) Err() error {
	return <-s.err
}
