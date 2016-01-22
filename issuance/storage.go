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
	"path/filepath"
	"strings"
)

// Workspace is where this program stores assets.
var Workspace = Storage(DefaultWorkspace)

// Storage is a root directory and facilitates forming file paths
// derived from it.
type Storage string

// Sites gets the directory that stores site certificate and keys.
func (s Storage) Sites() string {
	return filepath.Join(string(s), "sites")
}

// Site returns the path to the folder containing assets for domain.
func (s Storage) Site(domain string) string {
	return filepath.Join(s.Sites(), strings.ToLower(domain))
}

// SiteCertFile returns the path to the certificate file for domain.
func (s Storage) SiteCertFile(domain string) string {
	return filepath.Join(s.Site(domain), strings.ToLower(domain)+".crt")
}

// SiteKeyFile returns the path to domain's private key file.
func (s Storage) SiteKeyFile(domain string) string {
	return filepath.Join(s.Site(domain), strings.ToLower(domain)+".key")
}

// SiteMetaFile returns the path to the domain's asset metadata file.
func (s Storage) SiteMetaFile(domain string) string {
	return filepath.Join(s.Site(domain), strings.ToLower(domain)+".json")
}

// Users gets the directory that stores account folders.
func (s Storage) Users() string {
	return filepath.Join(string(s), "users")
}

// User gets the account folder for the user with email.
func (s Storage) User(email string) string {
	if email == "" {
		email = emptyEmail
	}
	return filepath.Join(s.Users(), strings.ToLower(email))
}

// UserRegFile gets the path to the registration file for
// the user with the given email address.
func (s Storage) UserRegFile(email string) string {
	if email == "" {
		email = emptyEmail
	}
	fileName := emailUsername(email)
	if fileName == "" {
		fileName = "registration"
	}
	return filepath.Join(s.User(email), strings.ToLower(fileName)+".json")
}

// UserKeyFile gets the path to the private key file for
// the user with the given email address.
func (s Storage) UserKeyFile(email string) string {
	if email == "" {
		email = emptyEmail
	}
	fileName := emailUsername(email)
	if fileName == "" {
		fileName = "private"
	}
	return filepath.Join(s.User(email), strings.ToLower(fileName)+".key")
}

// emailUsername returns the username portion of an
// email address (part before '@') or the original
// input if it can't find the "@" symbol.
func emailUsername(email string) string {
	at := strings.Index(email, "@")
	if at == -1 {
		return email
	} else if at == 0 {
		return email[1:]
	}
	return email[:at]
}

// The name of the folder for accounts where the email
// address was not provided.
const emptyEmail = "default"
