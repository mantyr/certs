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

// Package issuance facilitates the issuance of certificates
// via the ACME protocol.
package issuance

// ServerURL is the URL to the ACME CA's directory. This must be
// set before obtaining certificates.
var ServerURL string

// Agree is whether the user agrees to the CA's service agreement.
// This need only be true if the user has not agreed before.
var Agree bool

// DefaultWorkspace is where assets will be stored if no custom
// Workspace variable is set by the importing package.
const DefaultWorkspace = "./certs_data"
