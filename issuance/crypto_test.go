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
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"os"
	"runtime"
	"testing"

	"github.com/xenolf/lego/acme"
)

func init() {
	rsaKeySize = 128 // make tests faster; small key size OK for testing
}

func TestSaveAndLoadRSAPrivateKey(t *testing.T) {
	keyFile := "test.key"
	defer os.Remove(keyFile)

	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		t.Fatal(err)
	}

	// test save
	err = saveRSAPrivateKey(privateKey, keyFile)
	if err != nil {
		t.Fatal("error saving private key:", err)
	}

	// it doesn't make sense to test file permission on windows
	if runtime.GOOS != "windows" {
		// get info of the key file
		info, err := os.Stat(keyFile)
		if err != nil {
			t.Fatal("error stating private key:", err)
		}
		// verify permission of key file is correct
		if info.Mode().Perm() != 0600 {
			t.Error("Expected key file to have permission 0600, but it wasn't")
		}
	}

	// test load
	loadedKey, err := loadRSAPrivateKey(keyFile)
	if err != nil {
		t.Error("error loading private key:", err)
	}

	// verify loaded key is correct
	if !rsaPrivateKeysSame(privateKey, loadedKey) {
		t.Error("Expected key bytes to be the same, but they weren't")
	}
}

// rsaPrivateKeysSame compares the bytes of a and b and returns true if they are the same.
func rsaPrivateKeysSame(a, b *rsa.PrivateKey) bool {
	return bytes.Equal(rsaPrivateKeyBytes(a), rsaPrivateKeyBytes(b))
}

// rsaPrivateKeyBytes returns the bytes of DER-encoded key.
func rsaPrivateKeyBytes(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

func TestSaveCertResource(t *testing.T) {
	Workspace = Storage("./certs_test_save")
	defer func() {
		err := os.RemoveAll(string(Workspace))
		if err != nil {
			t.Fatalf("Could not remove temporary storage directory (%s): %v", Workspace, err)
		}
	}()

	domain := "example.com"
	certContents := "certificate"
	keyContents := "private key"
	metaContents := `{
	"domain": "example.com",
	"certUrl": "https://example.com/cert",
	"certStableUrl": "https://example.com/cert/stable"
}`

	cert := acme.CertificateResource{
		Domain:        domain,
		CertURL:       "https://example.com/cert",
		CertStableURL: "https://example.com/cert/stable",
		PrivateKey:    []byte(keyContents),
		Certificate:   []byte(certContents),
	}

	err := saveCertResource(cert)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	certFile, err := ioutil.ReadFile(Workspace.SiteCertFile(domain))
	if err != nil {
		t.Errorf("Expected no error reading certificate file, got: %v", err)
	}
	if string(certFile) != certContents {
		t.Errorf("Expected certificate file to contain '%s', got '%s'", certContents, string(certFile))
	}

	keyFile, err := ioutil.ReadFile(Workspace.SiteKeyFile(domain))
	if err != nil {
		t.Errorf("Expected no error reading private key file, got: %v", err)
	}
	if string(keyFile) != keyContents {
		t.Errorf("Expected private key file to contain '%s', got '%s'", keyContents, string(keyFile))
	}

	metaFile, err := ioutil.ReadFile(Workspace.SiteMetaFile(domain))
	if err != nil {
		t.Errorf("Expected no error reading meta file, got: %v", err)
	}
	if string(metaFile) != metaContents {
		t.Errorf("Expected meta file to contain '%s', got '%s'", metaContents, string(metaFile))
	}
}

func TestExistingCertAndKey(t *testing.T) {
	Workspace = Storage("./le_test_existing")
	defer func() {
		err := os.RemoveAll(string(Workspace))
		if err != nil {
			t.Fatalf("Could not remove temporary storage directory (%s): %v", Workspace, err)
		}
	}()

	domain := "example.com"

	if existingCertAndKey(domain) {
		t.Errorf("Did NOT expect %v to have existing cert or key, but it did", domain)
	}

	err := saveCertResource(acme.CertificateResource{
		Domain:      domain,
		PrivateKey:  []byte("key"),
		Certificate: []byte("cert"),
	})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !existingCertAndKey(domain) {
		t.Errorf("Expected %v to have existing cert and key, but it did NOT", domain)
	}
}
