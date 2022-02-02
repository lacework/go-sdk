//
// Author:: Salim Afiune Maya (<afiune@lacework.net>)
// Copyright:: Copyright 2020, Lacework Inc.
// License:: Apache License, Version 2.0
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
//

// A Lacework component package to help facilitate the loading and execution of components
package lwcomponent

import (
	"crypto/sha256"
	_ "embed"
	"encoding/base64"
	"strings"

	"aead.dev/minisign"
	"github.com/pkg/errors"
)

// Verify will verify a two-level signature.
// TODO(pjm): Can we make the signature format be an argument?
func verifySignature(rootKey minisign.PublicKey, file, sig []byte) error {
	h := sha256.New()
	if _, err := h.Write(file); err != nil {
		return errors.New("unable to compute hash for component")
	}

	var parsedSig minisign.Signature
	if err := parsedSig.UnmarshalText(sig); err != nil {
		return errors.New("unable to parse signature")
	}

	sigParts := strings.Split(parsedSig.TrustedComment, ".")
	if len(sigParts) != 4 {
		return errors.New("invalid signature trusted comment")
	}

	parsedRootSig, err := base64.StdEncoding.DecodeString(sigParts[3])
	if err != nil {
		return errors.Wrap(err, "unable to parse root signature from trusted comment")
	}

	if !minisign.Verify(rootKey, []byte(sigParts[1]), parsedRootSig) {
		return errors.New("invalid root signature over signing key")
	}

	var signingKey minisign.PublicKey
	if err := signingKey.UnmarshalText([]byte(sigParts[1])); err != nil {
		return errors.Wrap(err, "unable to parse signing key")
	}

	if !minisign.Verify(signingKey, h.Sum(nil), sig) {
		return errors.New("invalid signature over component")
	}

	return nil
}
