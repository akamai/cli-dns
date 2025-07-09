// Copyright 2020. Akamai Technologies, Inc
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
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateRecordsets(c *cli.Context) error {

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	var (
		zonename   string
		outputPath string
		inputPath  string
	)

	// Validate zonename argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	zonename = c.Args().First()

	// Check if the zone is an ALIAS zone
	zoneResp, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{
		Zone: zonename,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to retrieve zone information for %s. Error: %s", zonename, err)), 1)
	}
	if strings.EqualFold(zoneResp.Type, "ALIAS") {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Zone %s is an ALIAS zone and does not have recordsets", zonename)), 1)
	}

	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	} else {
		return cli.NewExitError(color.RedString("Input file is required"), 1)
	}

	// Parse input JSON file
	data, err := os.ReadFile(filepath.FromSlash(inputPath))
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to read input file"), 1)
	}
	recordsets := &dns.RecordSets{}
	err = json.Unmarshal(data, recordsets)
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to parse json file content"), 1)
	}

	// Determine update mode (overwrite or update existing recordset)
	var recordsetWorkList []dns.RecordSet

	if c.IsSet("overwrite") && c.Bool("overwrite") {
		recordsets := &dns.RecordSets{}
		err = json.Unmarshal(data, recordsets)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to parse json file content"), 1)
		}
		recordsetWorkList = recordsets.RecordSets
	} else {
		fmt.Println("Retrieving Existing Recordsets ", "")
		resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{
			Zone: zonename,
			QueryArgs: &dns.RecordSetQueryArgs{
				ShowAll: true,
			},
		})
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset List retrieval failed. Error: %s", err.Error())), 1)
		}

		fmt.Println("Processing Updated Recordsets ", "")
		recordsetWorkList = resp.RecordSets

		// Merge changes from input file
		soaInSet := false
		soaIndex := 0

		for _, crs := range recordsets.RecordSets {
			for i, rs := range recordsetWorkList {
				if crs.Name == rs.Name && crs.Type == rs.Type {
					recordsetWorkList[i] = crs
					if crs.Type == "SOA" {
						soaInSet = true
					}
				} else if rs.Type == "SOA" {
					soaIndex = i
				}
			}
		}

		// Auto-increment SOA serial if not explicitly set
		if !soaInSet && (soaIndex > 0 || recordsetWorkList[soaIndex].Type == "SOA") {
			soavals := strings.Split(recordsetWorkList[soaIndex].Rdata[0], " ")
			v, _ := strconv.Atoi(soavals[2])
			soavals[2] = strconv.Itoa(v + 1)
			recordsetWorkList[soaIndex].Rdata[0] = strings.Join(soavals, " ")
		}
	}

	// Submit recordset updates
	fmt.Println("Updating Recordsets ", "")
	recordsets.RecordSets = recordsetWorkList
	err = dnsClient.UpdateRecordSets(ctx, dns.UpdateRecordSetsRequest{
		Zone:       zonename,
		RecordSets: &dns.RecordSets{RecordSets: recordsetWorkList},
		RecLock:    []bool{true},
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset update failed. Error: %s", err.Error())), 1)
	}

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	// Fetch full updated list
	fmt.Fprintln(os.Stderr, color.BlueString("Retrieving full recordsets list...\n"))
	resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{
		Zone: zonename,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset List retrieval failed. Error: %s", err.Error())), 1)
	}

	results := ""

	// Format output as JSON or table format
	if c.IsSet("json") && c.Bool("json") {
		rjson, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to display recordsets list"), 1)
		}
		results = string(rjson)
	} else {
		results = renderRecordsetListTable(zonename, resp.RecordSets)
	}

	// Write output to file or console
	if len(outputPath) > 1 {
		//fmt.Printf("Writing Output to %s ", outputPath)
		rlfHandle, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer rlfHandle.Close()
		_, err = rlfHandle.WriteString(string(results))
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to write zone list output to file"), 1)
		}
		rlfHandle.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
