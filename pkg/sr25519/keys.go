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

// NewPk decodes a public key from the passed byte slice.
func NewPk(_data []byte) (*schnorrkel.PublicKey, error) {
	var data [32]byte
	copy(data[:], _data)
	pk := new(schnorrkel.PublicKey)
	return pk, pk.Decode(data)
}

// NewPkFromHex accepts hex strings with and without 0x.
func NewPkFromHex(hex string) (*schnorrkel.PublicKey, error) {
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	return NewPk(data)
}

// NewPkFromRng returns a new public key by readind data from rng.
func NewPkFromRng(rng io.Reader) (*schnorrkel.PublicKey, error) {
	sk, err := NewSkFromRng(rng)
	if err != nil {
		return nil, err
	}
	return sk.Public(), nil
}

// NewSk decodes a secret key from the passed byte slice.
func NewSk(_data []byte) (*schnorrkel.MiniSecretKey, error) {
	var data [32]byte
	copy(data[:], _data)
	return schnorrkel.NewMiniSecretKeyFromRaw(data)
}

// NewSkFromHex returns a secret key with the passed string as key.
// Accepts hex strings with and without 0x.
func NewSkFromHex(hex string) (*schnorrkel.MiniSecretKey, error) {
	data, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	return NewSk(data)
}

// NewSkFromRng returns a new secret key by readind data from rng.
func NewSkFromRng(rng io.Reader) (*schnorrkel.MiniSecretKey, error) {
	var data [32]byte
	if _, err := rng.Read(data[:]); err != nil {
		return nil, err
	}
	return NewSk(data[:])
}
