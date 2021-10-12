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

package sr25519

import (
	"io"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// NewPK decodes a public key from the passed byte slice.
func NewPK(data []byte) (*schnorrkel.PublicKey, error) {
	var _data [32]byte
	copy(_data[:], data)
	pk := new(schnorrkel.PublicKey)
	return pk, pk.Decode(_data)
}

// NewPKFromHex returns a public key by decoding a hex string.
// Accepts hex strings with and without 0x.
func NewPKFromHex(hex string) (*schnorrkel.PublicKey, error) {
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	return NewPK(data)
}

// NewPKFromRng returns a new public key by readind data from rng.
func NewPKFromRng(rng io.Reader) (*schnorrkel.PublicKey, error) {
	sk, err := NewSKFromRng(rng)
	if err != nil {
		return nil, err
	}
	return sk.Public(), nil
}

// NewSK decodes a secret key from the passed byte slice.
func NewSK(data []byte) (*schnorrkel.MiniSecretKey, error) {
	var _data [32]byte
	copy(_data[:], data)
	return schnorrkel.NewMiniSecretKeyFromRaw(_data)
}

// NewSKFromHex returns a secret key by decoding a hex string.
// Accepts hex strings with and without 0x.
func NewSKFromHex(hex string) (*schnorrkel.MiniSecretKey, error) {
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	return NewSK(data)
}

// NewSKFromRng returns a new secret key by readind data from rng.
func NewSKFromRng(rng io.Reader) (*schnorrkel.MiniSecretKey, error) {
	var data [32]byte
	if _, err := rng.Read(data[:]); err != nil {
		return nil, err
	}
	return NewSK(data[:])
}
