// Copyright 2018. Akamai Technologies, Inc
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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateZone(c *cli.Context) error {

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Session initialization failed: %v", err), 1)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Validate zonename argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}
	zonename := c.Args().First()

	var (
		inputPath  string
		outputPath string
	)

	var fileData []byte
	if c.IsSet("file") {
		inputPath = filepath.FromSlash(c.String("file"))
		fileData, err = os.ReadFile(inputPath)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read input file: %v", err)), 1)
		}
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return cli.NewExitError(color.RedString("No input file or piped data provided"), 1)
		}
		fileData, err = io.ReadAll(os.Stdin)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read from STDIN: %v", err)), 1)
		}
	}

	if c.IsSet("output") {
		outputPath = filepath.FromSlash(c.String("output"))
	}

	// Handle master zone file upload if --dns flag is set
	if c.Bool("dns") {
		masterZoneFileData := string(fileData)

		const httpMaxBody = 10 * 1024 * 1024
		if len(masterZoneFileData) > httpMaxBody {
			return cli.NewExitError(color.RedString("Master Zone File size too large to process"), 1)
		}

		fmt.Println("Uploading Master Zone File ...")
		err = dnsClient.PostMasterZoneFile(ctx, dns.PostMasterZoneFileRequest{
			Zone:     zonename,
			FileData: masterZoneFileData,
		})
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Master Zone File upload failed: %v", err)), 1)
		}
		fmt.Println("Master Zone File uploaded successfully.")
		return nil
	}

	// Parse input JSON for recordsets
	inputRecordSets := &dns.RecordSets{}
	err = json.Unmarshal(fileData, inputRecordSets)
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to parse JSON input file: %v", err)), 1)
	}

	// Prepare recordset update list and handle SOA record
	var recordsetWorkList []dns.RecordSet
	soaInSet := false
	soaIndex := -1

	if c.Bool("overwrite") {
		recordsetWorkList = inputRecordSets.RecordSets
		for i, rs := range recordsetWorkList {
			if rs.Type == "SOA" {
				soaInSet = true
				soaIndex = i
				break
			}
		}
	} else {
		fmt.Println("Retrieving Existing Recordsets ...")
		existingResp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{
			Zone: zonename,
			QueryArgs: &dns.RecordSetQueryArgs{
				ShowAll: true,
			},
		})
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset list retrieval failed: %v", err)), 1)
		}

		recordsetWorkList = existingResp.RecordSets

		for i, rs := range recordsetWorkList {
			if rs.Type == "SOA" {
				soaIndex = i
				break
			}
		}

		for _, updatedRS := range inputRecordSets.RecordSets {
			found := false
			for i, existingRS := range recordsetWorkList {
				if updatedRS.Name == existingRS.Name && updatedRS.Type == existingRS.Type {
					recordsetWorkList[i] = updatedRS
					found = true
					break
				}
			}
			if !found {
				recordsetWorkList = append(recordsetWorkList, updatedRS)
			}
			if updatedRS.Type == "SOA" {
				soaInSet = true
			}
		}

		if !soaInSet && soaIndex >= 0 {
			soaRec := &recordsetWorkList[soaIndex]
			if len(soaRec.Rdata) > 0 {
				soavals := strings.Fields(soaRec.Rdata[0])
				if len(soavals) >= 3 {
					serial, err := strconv.Atoi(soavals[2])
					if err == nil {
						serial++
						soavals[2] = strconv.Itoa(serial)
						soaRec.Rdata[0] = strings.Join(soavals, " ")
					} else {
						fmt.Fprintf(os.Stderr, "Warning: failed to parse SOA serial: %v\n", err)
					}
				}
			}
		}
	}

	fmt.Println("Updating Recordsets ...")
	err = dnsClient.UpdateRecordSets(ctx, dns.UpdateRecordSetsRequest{
		Zone:       zonename,
		RecordSets: &dns.RecordSets{RecordSets: recordsetWorkList},
		RecLock:    nil,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset update failed: %v", err)), 1)
	}

	if c.Bool("suppress") {
		return nil
	}

	// Retrieve and display updated recordsets
	fmt.Println("Retrieving Updated Recordsets ...")
	resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{
		Zone: zonename,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to retrieve recordsets after update: %v", err)), 1)
	}

	var results string
	if c.Bool("json") {
		rjson, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to display recordsets list"), 1)
		}
		results = string(rjson)
	} else {
		results = renderRecordsetListTable(zonename, resp.RecordSets)
	}

	// Output results to file or console
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file: %v", err)), 1)
		}
		defer f.Close()
		_, err = f.WriteString(results)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to write zone output to file"), 1)
		}
		f.Sync()
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
