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

package substrate

import (
	"github.com/centrifuge/go-substrate-rpc-client/v3/rpc/state"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"perun.network/go-perun/log"
	pkgsync "perun.network/go-perun/pkg/sync"
)

type (
	// EventSource collects all events from the chain.
	// Can then be used to filter out events, eg. for a specific pallet.
	EventSource struct {
		*pkgsync.Closer
		log.Embedding

		future *state.StorageSubscription
		api    *Api

		events chan types.EventRecordsRaw
		err    chan error
	}

	// EventKey can be used to construct a types.StorageKey.
	EventKey struct {
		firstName, lastName string
		args                [][]byte
	}
)

// GSRPC buffers up to 20k events.
const ChanBuffSize = 1024

// SystemEventsKey is the key for all system events.
func SystemEventsKey() *EventKey {
	return &EventKey{"System", "Events", nil}
}

// SystemAccountKey is the key to query an account.
func SystemAccountKey(accountId types.AccountID) *EventKey {
	return &EventKey{"System", "Account", [][]byte{accountId[:]}}
}

// Key returns the EventKey as GSRPC key.
func (k *EventKey) Key(meta *types.Metadata) (types.StorageKey, error) {
	return types.CreateStorageKey(meta, k.firstName, k.lastName, k.args...)
}

// NewEventSource returns a new EventSource.
// It queries pastBlocks into the past to retrive old events and starts listening
// for the EventKeys that were passed.
func NewEventSource(api *Api, pastBlocks types.BlockNumber, _keys ...*EventKey) (*EventSource, error) {
	log := log.MakeEmbedding(log.Get())
	meta := api.Metadata()
	var keys []types.StorageKey
	for _, _key := range _keys {
		key, err := _key.Key(meta)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	// Subscribe to new events
	future, err := api.Subscribe(keys...)
	if err != nil {
		return nil, err
	}

	events := make(chan types.EventRecordsRaw, ChanBuffSize)
	errChan := make(chan error, 16)
	ret := &EventSource{Closer: new(pkgsync.Closer), Embedding: log, future: future, api: api, events: events, err: errChan}
	ret.OnCloseAlways(func() {
		future.Unsubscribe()
	})
	return ret, ret.init(keys, meta, pastBlocks)
}

func (s *EventSource) init(keys []types.StorageKey, meta *types.Metadata, pastBlocks types.BlockNumber) error {
	header, err := s.api.LastHeader()
	if err != nil {
		return err
	}
	startNum := types.BlockNumber(0)
	if header.Number > pastBlocks {
		startNum = header.Number - pastBlocks
	}
	s.Log().WithField("startNum", startNum).Tracef("Calculated start block num, past=%v", pastBlocks)
	startBlock, err := s.api.BlockHash(uint64(startNum))
	if err != nil {
		return err
	}
	// Query and parse past events
	past, err := s.api.QueryAll(keys, startBlock)
	if err != nil {
		return err
	}
	s.Log().WithField("pastEvents", len(past)).Debug("PalletEventSource started")
	// This could dead-lock if there are more events than the channel can hold.
	for _, set := range past {
		s.parseEvent(set, meta)
	}

	go func() {
		defer s.Log().Debug("PalletEventSource stopped")
		defer close(s.err)
		defer s.Close()

		for {
			select {
			case set := <-s.future.Chan():
				s.parseEvent(set, meta)
			case err := <-s.future.Err():
				s.err <- err
				return
			case <-s.Closed():
				return
			}
		}
	}()
	return nil
}

func (s *EventSource) parseEvent(set types.StorageChangeSet, meta *types.Metadata) {
	for _, change := range set.Changes {
		if !change.HasStorageData {
			continue
		}
		record := types.EventRecordsRaw(change.StorageData)
		s.events <- record
	}
}

// Events returns a channel that contains all events that the EventSource found.
// This channel will never be closed, use Err() instead.
func (s *EventSource) Events() <-chan types.EventRecordsRaw {
	return s.events
}

// Err returns the error channel. Will be closed when the EventSource gets closed.
func (s *EventSource) Err() <-chan error {
	return s.err
}
