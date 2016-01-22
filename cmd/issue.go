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

package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mholt/certs/issuance"
	"github.com/spf13/cobra"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Issue new certificates in bulk",
	Long: `The issue command will obtain new certificates according to an
input file. The input file may be in text or CSV format.

Each line in the input file will issue a certificate. Multiple
domains may appear on a single line, separated by the delimiter.
The first domain on each line will be the Common Name on the
certificate, and all others will be SubjectAltName entries.

Certificates are stored with the associated private key in a
folder in the workspace (customized with --out) in a subfolder
named after the Common Name on the certificate. If a cert with
the same Common Name already exists in the workspace, that
certificate will be skipped, even if its SAN entries differ.`,
	Run: runIssue,
}

func runIssue(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("missing argument: input file with list of domains")
	}

	workspaceDir, err := cmd.Flags().GetString("out")
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	issuance.Workspace = issuance.Storage(workspaceDir)

	delim, err := cmd.Flags().GetString("delim")
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	domainList, domainMap, err := loadDomains(args[0], delim)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	// addwww, err := cmd.Flags().GetString("addwww")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// addwww = strings.ToLower(addwww)
	// if addwww != "separate" && addwww != "san" && addwww != "cn" {
	// 	log.Fatalf("addwww: '%s' is not recognized; valid values are CN, SAN, or seperate", addwww)
	// }

	log.Printf("[INFO] Obtaining %d certificates for %d domains\n", len(domainList), len(domainMap))

	ca, err := cmd.Flags().GetString("ca")
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	issuance.ServerURL = ca

	agree, err := cmd.Flags().GetBool("agree")
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	issuance.Agree = agree

	email, err := cmd.Flags().GetString("email")
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	user, err := issuance.GetUser(email)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	if err := user.ObtainCerts(domainList); err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
}

// loadDomains loads a list of domain names from filename, separated
// and returns each domain in a slice (ordered) and in a map (unordered).
// These dual return values allow you to iterate in order, but also do
// quick membership checking, despite using a more memory.
func loadDomains(filename string, delim string) ([][]string, map[string]struct{}, error) {
	comma, err := standardDelimiter(delim)
	if err != nil {
		return nil, nil, err
	}

	inputFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("loading input file: %v", err)
	}
	defer inputFile.Close()

	var domainList [][]string
	var domainMap = make(map[string]struct{})

	r := csv.NewReader(inputFile)
	r.Comma = comma
	r.FieldsPerRecord = -1

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return domainList, domainMap, err
		}

		var cleanedRow []string
		for _, val := range row {
			domain := strings.TrimSpace(strings.ToLower(val))
			if domain != "" {
				cleanedRow = append(cleanedRow, domain)
				domainMap[domain] = struct{}{}
			}
		}
		if len(cleanedRow) > 0 {
			domainList = append(domainList, cleanedRow)
		}
	}

	return domainList, domainMap, nil
}

func standardDelimiter(delim string) (rune, error) {
	delim = strings.ToLower(delim)
	if delim == "" {
		delim = defaultDelimiter
	}
	if delim == "." || delim == "\n" || delim == "\r" || delim == "\"" || delim == "'" {
		return 0, fmt.Errorf("'%s' is not a valid delimiter for this program", delim)
	}
	if delim == "\\t" || delim == "tab" {
		return '\t', nil
	}
	if delim == "asc30" { // record separator
		return 30, nil
	}
	if delim == "asc31" { // unit separator
		return 31, nil
	}
	if len(delim) != 1 {
		return 0, fmt.Errorf("'%s' is not a valid delimiter; can only be 1 character", delim)
	}
	return rune(delim[0]), nil
}

func init() {
	RootCmd.AddCommand(issueCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issueCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//issueCmd.Flags().String("addwww", "", "Ensure www variant is added for every domain")
	issueCmd.Flags().String("delim", defaultDelimiter, "Delimiter")
	issueCmd.Flags().String("ca", "https://acme-staging.api.letsencrypt.org/directory", "URL of directory for ACME server")
	issueCmd.Flags().String("email", "", "Email address to register with CA for account recovery")
	issueCmd.Flags().String("out", issuance.DefaultWorkspace, "Path to folder in which to store assets")
	issueCmd.Flags().Bool("agree", false, "Indicate your agreement to CA's legal terms")
}

const defaultDelimiter = ","
