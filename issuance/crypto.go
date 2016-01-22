// Copyright © 2016 Matthew Holt
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

package issuance

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"os"

	"github.com/xenolf/lego/acme"
)

// loadRSAPrivateKey loads a PEM-encoded RSA private key from file.
func loadRSAPrivateKey(file string) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyBytes)
	return x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
}

// saveRSAPrivateKey saves a PEM-encoded RSA private key to file.
func saveRSAPrivateKey(key *rsa.PrivateKey, file string) error {
	pemKey := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}
	keyOut, err := os.Create(file)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	keyOut.Chmod(0600)
	return pem.Encode(keyOut, &pemKey)
}

// saveCertResource saves the certificate resource to disk. This
// includes the certificate file itself, the private key, and the
// metadata file.
func saveCertResource(cert acme.CertificateResource) error {
	err := os.MkdirAll(Workspace.Site(cert.Domain), 0700)
	if err != nil {
		return err
	}

	// Save cert
	err = ioutil.WriteFile(Workspace.SiteCertFile(cert.Domain), cert.Certificate, 0600)
	if err != nil {
		return err
	}

	// Save private key
	err = ioutil.WriteFile(Workspace.SiteKeyFile(cert.Domain), cert.PrivateKey, 0600)
	if err != nil {
		return err
	}

	// Save cert metadata
	jsonBytes, err := json.MarshalIndent(&cert, "", "\t")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(Workspace.SiteMetaFile(cert.Domain), jsonBytes, 0600)
	if err != nil {
		return err
	}

	return nil
}

// existingCertAndKey returns true if the host has a certificate
// and private key in storage already, false otherwise.
func existingCertAndKey(host string) bool {
	_, err := os.Stat(Workspace.SiteCertFile(host))
	if err != nil {
		return false
	}
	_, err = os.Stat(Workspace.SiteKeyFile(host))
	if err != nil {
		return false
	}
	return true
}

// rsaKeySize is the size in bits to use for new RSA keys.
// This shouldn't need to change except for in tests; the
// size can be drastically reduced for speed.
var rsaKeySize = 2048
