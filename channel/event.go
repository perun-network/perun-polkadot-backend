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
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
)

type (
	// EventRecords contains all events that can be emitted by a substrate chain
	// that contains the perun pallet.
	EventRecords struct {
		types.EventRecords

		PerunModule_Deposited []DepositedEvent
		PerunModule_Disputed  []DisputedEvent
		PerunModule_Concluded []ConcludedEvent
		PerunModule_Withdrawn []WithdrawnEvent
	}

	// PerunEvent is a perun event.
	PerunEvent interface{}

	// DepositedEvent is emitted when a funding id receives a deposit.
	DepositedEvent struct {
		Phase   types.Phase // required
		Fid     FundingId
		Balance Balance
		Topics  []types.Hash // required
	}

	// DisputedEvent is emitted when a dispute was opened or updated.
	DisputedEvent struct {
		Phase  types.Phase // required
		Cid    ChannelId
		State  State
		Topics []types.Hash // required
	}

	// ConcludedEvent is emitted when a channel is concluded.
	ConcludedEvent struct {
		Phase  types.Phase // required
		Cid    ChannelId
		Topics []types.Hash // required
	}

	// WithdrawnEvent is emitted when all funds are withdrawn from a funding id.
	WithdrawnEvent struct {
		Phase  types.Phase // required
		Fid    FundingId
		Topics []types.Hash // required
	}
)

// EventIsPerunEvent returns whether an event is a PerunEvent.
func EventIsPerunEvent(e interface{}) bool {
	switch e.(type) {
	case *DepositedEvent:
		return true // fallthrough is not possible in type switches
	case *DisputedEvent:
		return true
	case *ConcludedEvent:
		return true
	case *WithdrawnEvent:
		return true
	default:
		return false
	}
}

// EventIsDeposited returns whether an event is a DepositedEvent.
func EventIsDeposited(e PerunEvent) bool {
	_, ok := e.(*DepositedEvent)
	return ok
}

// EventIsDisputed returns whether an event is a DisputedEvent.
func EventIsDisputed(cid ChannelId) func(PerunEvent) bool {
	return func(e PerunEvent) bool {
		event, ok := e.(*DisputedEvent)
		return ok && event.Cid == cid
	}
}

// EventIsConcluded returns whether an event is a ConcludedEvent.
func EventIsConcluded(cid ChannelId) func(PerunEvent) bool {
	return func(e PerunEvent) bool {
		event, ok := e.(*ConcludedEvent)
		return ok && event.Cid == cid
	}
}

// EventIsWithdrawn returns whether an event is a WithdrawnEvent.
func EventIsWithdrawn(e PerunEvent) bool {
	_, ok := e.(*WithdrawnEvent)
	return ok
}

// Merge merges all perun events into one slice and returns it.
func (r *EventRecords) Merge() []PerunEvent {
	var ret []PerunEvent
	// The order is not changed here since all but one slice are empty.
	for _, e := range r.PerunModule_Deposited {
		e := e
		ret = append(ret, &e)
	}
	for _, e := range r.PerunModule_Disputed {
		e := e
		ret = append(ret, &e)
	}
	for _, e := range r.PerunModule_Concluded {
		e := e
		ret = append(ret, &e)
	}
	for _, e := range r.PerunModule_Withdrawn {
		e := e
		ret = append(ret, &e)
	}
	return ret
}
