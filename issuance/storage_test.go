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
	"testing"
)

func TestStorage(t *testing.T) {
	Workspace = Storage("./certs_test")

	if expected, actual := filepath.Join("certs_test", "sites"), Workspace.Sites(); actual != expected {
		t.Errorf("Expected Sites() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "sites", "test.com"), Workspace.Site("Test.com"); actual != expected {
		t.Errorf("Expected Site() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "sites", "test.com", "test.com.crt"), Workspace.SiteCertFile("Test.com"); actual != expected {
		t.Errorf("Expected SiteCertFile() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "sites", "test.com", "test.com.key"), Workspace.SiteKeyFile("Test.com"); actual != expected {
		t.Errorf("Expected SiteKeyFile() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "sites", "test.com", "test.com.json"), Workspace.SiteMetaFile("Test.com"); actual != expected {
		t.Errorf("Expected SiteMetaFile() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "users"), Workspace.Users(); actual != expected {
		t.Errorf("Expected Users() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "users", "me@example.com"), Workspace.User("Me@example.com"); actual != expected {
		t.Errorf("Expected User() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "users", "me@example.com", "me.json"), Workspace.UserRegFile("Me@example.com"); actual != expected {
		t.Errorf("Expected UserRegFile() to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "users", "me@example.com", "me.key"), Workspace.UserKeyFile("Me@example.com"); actual != expected {
		t.Errorf("Expected UserKeyFile() to return '%s' but got '%s'", expected, actual)
	}

	// Test with empty emails
	if expected, actual := filepath.Join("certs_test", "users", emptyEmail), Workspace.User(emptyEmail); actual != expected {
		t.Errorf("Expected User(\"\") to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "users", emptyEmail, emptyEmail+".json"), Workspace.UserRegFile(""); actual != expected {
		t.Errorf("Expected UserRegFile(\"\") to return '%s' but got '%s'", expected, actual)
	}
	if expected, actual := filepath.Join("certs_test", "users", emptyEmail, emptyEmail+".key"), Workspace.UserKeyFile(""); actual != expected {
		t.Errorf("Expected UserKeyFile(\"\") to return '%s' but got '%s'", expected, actual)
	}
}

func TestEmailUsername(t *testing.T) {
	for i, test := range []struct {
		input, expect string
	}{
		{
			input:  "username@example.com",
			expect: "username",
		},
		{
			input:  "plus+addressing@example.com",
			expect: "plus+addressing",
		},
		{
			input:  "me+plus-addressing@example.com",
			expect: "me+plus-addressing",
		},
		{
			input:  "not-an-email",
			expect: "not-an-email",
		},
		{
			input:  "@foobar.com",
			expect: "foobar.com",
		},
		{
			input:  emptyEmail,
			expect: emptyEmail,
		},
		{
			input:  "",
			expect: "",
		},
	} {
		if actual := emailUsername(test.input); actual != test.expect {
			t.Errorf("Test %d: Expected username to be '%s' but was '%s'", i, test.expect, actual)
		}
	}
}
