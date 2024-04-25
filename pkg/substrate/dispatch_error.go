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
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
)

var (
	// ErrUnknownError the error could not be decoded.
	ErrUnknownError = errors.New("unknown dispatch error type")
	// ErrCallFailed an extrinsic returned an error.
	ErrCallFailed = errors.New("call failed")
)

// DecodeError decodes an error into a human-readable form.
// Returns either ErrUnknownError or ErrCallFailed.
func DecodeError(meta *types.Metadata, err types.DispatchError) error {
	metaV, ok := Meta(meta)
	if !ok {
		return errors.Wrap(ErrUnknownError, "wrong meta data version")
	}
	if err.IsToken {
		return errors.Wrap(ErrCallFailed, "token error")

	}
	if !err.IsModule {
		return errors.Wrap(ErrUnknownError, "not module error")
	}
	metaDataErr, findErr := metaV.FindError(err.ModuleError.Index, err.ModuleError.Error)
	if findErr != nil {
		return errors.Wrap(ErrUnknownError, findErr.Error())
	}
	return errors.Wrap(ErrCallFailed, formatErrorMeta(metaDataErr))
}

// formatErrorMeta formats an MetadataError into a human-readable form.
func formatErrorMeta(err *types.MetadataError) string {

	return fmt.Sprintf("%s:%s", err.Name, err.Value)
}
