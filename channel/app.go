// Copyright 2024 - See NOTICE file for copyright holders.
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
	"math/rand"

	dotwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519"
	dottestwallet "github.com/perun-network/perun-polkadot-backend/wallet/sr25519/test"
	"perun.network/go-perun/channel"
)

var _ channel.AppID = new(AppID)

// AppID is an app identifier and
// OffIdentity is the type of the app identifier (serialized address).
type AppID struct {
	OffIdentity
}

// AppIDKey is the key representation of an app identifier.
type AppIDKey string

// Equal compares two AppID objects for equality.
func (a AppID) Equal(b channel.AppID) bool {
	bTyped, ok := b.(*AppID)
	if !ok {
		return false
	}
	return a.OffIdentity == bTyped.OffIdentity

}

// Key returns the key representation of this app identifier.
func (a AppID) Key() channel.AppIDKey {
	b, err := a.MarshalBinary()
	if err != nil {
		panic(err)
	}
	return channel.AppIDKey(b)
}

// MarshalBinary marshals the contents of AppID into a byte string.
func (a AppID) MarshalBinary() ([]byte, error) {
	data := a.OffIdentity
	return data[:], nil
	// return data, nil
}

// UnmarshalBinary converts a bytestring, representing AppID into the AppID struct.
func (a *AppID) UnmarshalBinary(data []byte) error {
	addr := &dotwallet.Address{}
	err := addr.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	var appIdent [OffIdentityLen]byte
	copy(appIdent[:], data)

	appaddr := &AppID{appIdent}
	*a = *appaddr
	return nil
}

func NewRandomAppID(rng *rand.Rand) *AppID {
	addr := dottestwallet.NewRandomAddress(rng)

	// Assuming addr is of type wallet.Address, and we need sr25519.Address
	// Convert or cast the address to sr25519.Address if possible
	srAddr, ok := addr.(*dotwallet.Address)
	if !ok {
		panic("address is not of type *sr25519.Address")
	}
	appIdent, err := srAddr.MarshalBinary()
	if err != nil {
		panic(err)
	}

	var offIdentity OffIdentity
	copy(offIdentity[:], appIdent)

	return &AppID{offIdentity}
}
