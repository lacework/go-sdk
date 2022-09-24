//
// Author:: Nicholas Schmeller (<nick.schmeller@lacework.net>)
// Copyright:: Copyright 2022, Lacework Inc.
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

package lwrunner

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"golang.org/x/crypto/ssh"
)

// Helper function to generate a random private key. This key will
// stay in memory and does not persist across test runs.
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// generate key
	priv, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// validate key
	err = priv.Validate()
	if err != nil {
		return nil, err
	}

	return priv, err
}

// Generates SSH keys in EC2-readable format
// Returns public key bytes, private key bytes, error
func GetKeyBytes() ([]byte, []byte, error) {
	// generate keys with bit size 4096
	priv, err := GeneratePrivateKey(4096)
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := GetPublicKeyBytes(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	privBytes, err := EncodePrivateKeyToPEM(priv)
	if err != nil {
		return nil, nil, err
	}

	return pubBytes, privBytes, err
}

// Takes an rsa.PublicKey and returns bytes suitable for writing to .pub file
// Returns in the format "ssh-rsa ..."
func GetPublicKeyBytes(key *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(key)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, err
}

// Encodes private key from RSA struct to PEM formatted bytes
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) ([]byte, error) {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)
	if privatePEM == nil {
		return nil, errors.New("failed to encode private key to PEM")
	}

	return privatePEM, nil
}
