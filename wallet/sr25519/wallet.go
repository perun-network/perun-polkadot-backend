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
	"sync"

	pwallet "perun.network/go-perun/wallet"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/centrifuge/go-substrate-rpc-client/v3/types"
	"github.com/perun-network/perun-polkadot-backend/pkg/substrate"
	"github.com/pkg/errors"
)

// Wallet implements the Wallet interface. It uses sr25519 cryptography.
type Wallet struct {
	mtx     *sync.RWMutex               // protects all
	keyRing map[pwallet.AddrKey]keyPair // address -> keyPair
}

var (
	// WrongAddrType is returned when the type of a perun Address was not
	// actually of type Address.
	WrongAddrType = errors.New("got wrong address type")
	// AccountNotPresent is returned when an Account for a specific Address was
	// not found in the wallet.
	AccountNotPresent = errors.New("account is not present in the wallet")
)

// NewWallet returns a new wallet.
func NewWallet() *Wallet {
	return &Wallet{
		mtx:     new(sync.RWMutex),
		keyRing: make(map[pwallet.AddrKey]keyPair),
	}
}

// ImportSK can be used to import a secret key into the wallet.
func (w *Wallet) ImportSK(sk *schnorrkel.MiniSecretKey) *Account {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	// Create the account.
	keyPair := makeKeyPair(sk)
	// Insert and overwrite existing.
	w.keyRing[keyPair.pk.key()] = keyPair
	return newAccount(w, keyPair.pk)
}

// Unlock unlocks an Account. Needed by the perun Wallet interface.
// Returns WrongAddrType or AccountNotPresent in case of an error.
func (w *Wallet) Unlock(pAddr pwallet.Address) (pwallet.Account, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	if !IsAddr(pAddr) {
		return nil, WrongAddrType
	}
	addr := AsAddr(pAddr)
	key := pwallet.AddrKey(addr.key())
	// Check if the account exists.
	if _, found := w.keyRing[key]; !found {
		return nil, AccountNotPresent
	}
	return newAccount(w, addr), nil
}

// SignTx can sign Extrinsics for an address that is present in the wallet.
func (w *Wallet) signExt(addr *Address, tx *types.Extrinsic, opts types.SignatureOptions, netId substrate.NetworkId) error {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	kp, found := w.keyRing[addr.key()]
	if !found {
		return AccountNotPresent
	}
	// Temporarily derive a KeyringPair for signing. GSRPC requires it.
	gsrpcKp, err := kp.keyRing(netId)
	if err != nil {
		return errors.WithMessage(err, "deriving keypair")
	}
	return tx.Sign(gsrpcKp, opts)
}

// signData can sign arbitrary data for an address that is present in the wallet.
func (w *Wallet) signData(addr *Address, data []byte) ([]byte, error) {
	w.mtx.RLock()
	defer w.mtx.RUnlock()

	kp, found := w.keyRing[addr.key()]
	if !found {
		return nil, AccountNotPresent
	}
	// Actual signing.
	ctx := schnorrkel.NewSigningContext(SignaturePrefix, data)
	_sig, err := kp.sk.Sign(ctx)
	if err != nil {
		return nil, err
	}
	sig := _sig.Encode()
	return sig[:], nil
}

// LockAll does nothing. Needed by the perun Wallet interface.
func (*Wallet) LockAll() {}

// IncrementUsage does nothing. Needed by the perun Wallet interface.
func (*Wallet) IncrementUsage(pwallet.Address) {}

// DecrementUsage does nothing. Needed by the perun Wallet interface.
func (*Wallet) DecrementUsage(pwallet.Address) {}
