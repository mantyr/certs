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
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/xenolf/lego/acme"
)

// User is type that can interact with an ACME server.
type User struct {
	rateLimiter  `json:"-"`
	Email        string
	Registration *acme.RegistrationResource
	key          *rsa.PrivateKey
}

// GetUser loads the user with the given email from disk.
// If the user does not exist, it will create a new one,
// but it will NOT save new user to the disk or register
// it via ACME.
func GetUser(email string) (*User, error) {
	var user User

	// open user file
	regFile, err := os.Open(Workspace.UserRegFile(email))
	if err != nil {
		if os.IsNotExist(err) {
			// create a new user
			return newUser(email)
		}
		return nil, err
	}
	defer regFile.Close()

	// load user information
	err = json.NewDecoder(regFile).Decode(&user)
	if err != nil {
		return nil, err
	}

	// load their private key
	user.key, err = loadRSAPrivateKey(Workspace.UserKeyFile(email))
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// newUser creates a new User for the given email address
// with a new private key. This function will NOT save the
// user to disk or register it via ACME. If you want to use
// a user account that might already exist, call getUser
// instead.
func newUser(email string) (*User, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaKeySize)
	if err != nil {
		return nil, err
	}

	u := &User{
		Email: email,
		key:   privateKey,
	}

	return u, nil
}

// saveUser persists a user's key and account registration
// to the file system. It does NOT register the user via ACME.
func saveUser(user *User) error {
	// make user account folder
	err := os.MkdirAll(Workspace.User(user.Email), 0700)
	if err != nil {
		return err
	}

	// save private key file
	err = saveRSAPrivateKey(user.key, Workspace.UserKeyFile(user.Email))
	if err != nil {
		return err
	}

	// save registration file
	jsonBytes, err := json.MarshalIndent(user, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(Workspace.UserRegFile(user.Email), jsonBytes, 0600)
}

// ObtainCerts obtains certificates in bundles, where each slice in the slice
// is a list of domains to put onto the certificate. This function is robust
// in handling rate limiting and will retry until it succeeds.
func (u *User) ObtainCerts(bundles [][]string) error {
	if ServerURL == "" {
		return fmt.Errorf("must set ServerURL before obtaining certificates")
	}

	client, err := u.newClient()
	if err != nil {
		return err
	}

	for _, domains := range bundles {
		if len(domains) == 0 {
			log.Println("[INFO] Skipping a bundle with no domains specified")
			continue
		}

	Obtain:
		// certificate and key could have appeared since we last checked, especially if waiting for rate limit
		if existingCertAndKey(domains[0]) {
			log.Printf("[INFO] Existing certificate and key for %s. Skipping bundle: %v", domains[0], domains)
			continue
		}

		certRes, failures := client.ObtainCertificate(domains, true, nil)
		if len(failures) > 0 {
			for domain, err := range failures {
				if strings.Contains(err.Error(), "rateLimited") {
					u.BackOff()
					log.Printf("[WARNING][%s] Rate limited: %v - backing off and retrying in %v", domain, err, u.interval)
					u.Wait()
					log.Printf("Retrying certificate for %v", domains)
					goto Obtain
				} else if _, ok := err.(acme.TOSError); ok {
					log.Printf("[WARNING][%s] Updated legal terms: %v", domain, err)
					err := client.AgreeToTOS()
					if err != nil {
						return fmt.Errorf("error agreeing to updated terms for %s: %v", domain, err)
					}
					goto Obtain
				}
			}
			return ObtainError(failures)
		}

		// immediately save each certificate as we obtain it
		err := saveCertResource(certRes)
		if err != nil {
			return fmt.Errorf("error saving assets for %v: %v", domains, err)
		}

		// open throttle if it wasn't already
		u.Resume()
	}

	return nil
}

// newClient makes a new ACME client for the user u, including
// registering the user, agreeing to terms, and saving the user
// data to storage if the user was not already registered. The
// returned acme.Client is ready to use.
func (u *User) newClient() (*acme.Client, error) {
	client, err := acme.NewClient(ServerURL, u, rsaKeySize)
	if err != nil {
		return nil, fmt.Errorf("creating ACME client: %v", err)
	}

	// TODO: Customize ports and challenges

	if u.Registration == nil {
		if !Agree {
			return nil, fmt.Errorf("cannot register user '%s' without --agree", u.Email)
		}

		reg, err := client.Register()
		if err != nil {
			return nil, fmt.Errorf("registration error: %v", err)
		}
		u.Registration = reg

		err = client.AgreeToTOS()
		if err != nil {
			return nil, fmt.Errorf("error agreeing to terms: %v", err)
		}

		err = saveUser(u)
		if err != nil {
			return nil, fmt.Errorf("could not save user: %v", err)
		}
	}

	return client, nil
}

// GetEmail gets u's email.
func (u *User) GetEmail() string {
	return u.Email
}

// GetRegistration gets u's registration resource.
func (u *User) GetRegistration() *acme.RegistrationResource {
	return u.Registration
}

// GetPrivateKey gets u's private key.
func (u *User) GetPrivateKey() *rsa.PrivateKey {
	return u.key
}

// ObtainError maps failures keyed by domain
// name to their error message.
type ObtainError map[string]error

// Error returns a formatted, descriptive error message of failures in e.
func (e ObtainError) Error() string {
	var errMsg string
	for domain, err := range e {
		errMsg += "[" + domain + "] failed to get certificate: " + err.Error() + "\n"
	}
	return errMsg
}
