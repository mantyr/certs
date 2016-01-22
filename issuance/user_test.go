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
	"os"
	"testing"

	"github.com/xenolf/lego/acme"
)

func TestUser(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 128)
	if err != nil {
		t.Fatalf("Could not generate test private key: %v", err)
	}
	u := User{
		Email:        "me@mine.com",
		Registration: new(acme.RegistrationResource),
		key:          privateKey,
	}

	if expected, actual := "me@mine.com", u.GetEmail(); actual != expected {
		t.Errorf("Expected email '%s' but got '%s'", expected, actual)
	}
	if u.GetRegistration() == nil {
		t.Error("Expected a registration resource, but got nil")
	}
	if expected, actual := privateKey, u.GetPrivateKey(); actual != expected {
		t.Errorf("Expected the private key at address %p but got one at %p instead ", expected, actual)
	}
}

func TestNewUser(t *testing.T) {
	email := "me@foobar.com"
	user, err := newUser(email)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}
	if user.key == nil {
		t.Error("Private key is nil")
	}
	if user.Email != email {
		t.Errorf("Expected email to be %s, but was %s", email, user.Email)
	}
	if user.Registration != nil {
		t.Error("New user already has a registration resource; it shouldn't")
	}
}

func TestSaveUser(t *testing.T) {
	Workspace = Storage("./testdata")
	defer os.RemoveAll(string(Workspace))

	email := "me@foobar.com"
	user, err := newUser(email)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	err = saveUser(user)
	if err != nil {
		t.Fatalf("Error saving user: %v", err)
	}
	_, err = os.Stat(Workspace.UserRegFile(email))
	if err != nil {
		t.Errorf("Cannot access user registration file, error: %v", err)
	}
	_, err = os.Stat(Workspace.UserKeyFile(email))
	if err != nil {
		t.Errorf("Cannot access user private key file, error: %v", err)
	}
}

func TestGetUserDoesNotAlreadyExist(t *testing.T) {
	Workspace = Storage("./testdata")
	defer os.RemoveAll(string(Workspace))

	user, err := GetUser("user_does_not_exist@foobar.com")
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	if user.key == nil {
		t.Error("Expected user to have a private key, but it was nil")
	}
}

func TestGetUserAlreadyExists(t *testing.T) {
	Workspace = Storage("./testdata")
	defer os.RemoveAll(string(Workspace))

	email := "me@foobar.com"

	// Set up test
	user, err := newUser(email)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}
	err = saveUser(user)
	if err != nil {
		t.Fatalf("Error saving user: %v", err)
	}

	// Expect to load user from disk
	user2, err := GetUser(email)
	if err != nil {
		t.Fatalf("Error getting user: %v", err)
	}

	// Assert keys are the same
	if !rsaPrivateKeysSame(user.key, user2.key) {
		t.Error("Expected private key to be the same after loading, but it wasn't")
	}

	// Assert emails are the same
	if user.Email != user2.Email {
		t.Errorf("Expected emails to be equal, but was '%s' before and '%s' after loading", user.Email, user2.Email)
	}
}
