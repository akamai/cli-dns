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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateRecordsets(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing recordsets", "session failed %v", err)
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
		return failStep("Preparing recordsets", "zonename is required")
	}

	zonename = c.Args().First()

	// Check if the zone is an ALIAS zone
	zoneResp, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{
		Zone: zonename,
	})
	if err != nil {
		return failStep("Preparing recordsets", "Failed to retrieve zone information for %s. Error: %s", zonename, err)
	}
	if strings.EqualFold(zoneResp.Type, "ALIAS") {
		return failStep("Preparing recordsets", "Zone %s is an ALIAS zone and does not have recordsets", zonename)
	}

	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	} else {
		return failStep("Preparing recordsets", "Input file is required")
	}

	// Parse input JSON file
	data, err := os.ReadFile(filepath.FromSlash(inputPath))
	if err != nil {
		return failStep("Preparing recordsets", "Failed to read input file")
	}
	recordsets := &dns.RecordSets{}
	err = json.Unmarshal(data, recordsets)
	if err != nil {
		return failStep("Preparing recordsets", "Failed to parse json file content")
	}

	// Determine update mode (overwrite or update existing recordset)
	var recordsetWorkList []dns.RecordSet

	if c.IsSet("overwrite") && c.Bool("overwrite") {
		recordsets := &dns.RecordSets{}
		err = json.Unmarshal(data, recordsets)
		if err != nil {
			return failStep("Preparing recordsets", "Failed to parse json file content")
		}
		recordsetWorkList = recordsets.RecordSets
		fmt.Printf("Preparing recordsets ... %s\n", color.GreenString("[OK]"))
	} else {

		resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{
			Zone: zonename,
			QueryArgs: &dns.RecordSetQueryArgs{
				ShowAll: true,
			},
		})
		if err != nil {
			return failStep("Retrieving Existing Recordsets", "Recordset List retrieval failed. Error: %s", err.Error())
		}
		fmt.Fprintf(os.Stderr, "Retrieving Existing Recordsets ... %s\n", color.GreenString("[OK]"))
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
		fmt.Fprintf(os.Stderr, "Processing Updated Recordsets ... %s\n", color.GreenString("[OK]"))
	}

	// Submit recordset updates

	recordsets.RecordSets = recordsetWorkList
	err = dnsClient.UpdateRecordSets(ctx, dns.UpdateRecordSetsRequest{
		Zone:       zonename,
		RecordSets: &dns.RecordSets{RecordSets: recordsetWorkList},
		RecLock:    []bool{true},
	})
	if err != nil {
		return failStep("Updating Recordsets", "Recordset update failed. Error: %s", err.Error())
	}

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	fmt.Printf("Updating Recordsets ... %s\n", color.GreenString("[OK]"))
	// Fetch full updated list

	resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{
		Zone: zonename,
	})
	if err != nil {
		return failStep("Retrieving Recordsets List", "Recordset List retrieval failed. Error: %s", err.Error())
	}
	fmt.Fprintf(os.Stderr, "Retrieving Recordsets List ... %s\n", color.GreenString("[OK]"))
	results := ""

	// Format output as JSON or table format
	if c.IsSet("json") && c.Bool("json") {
		rjson, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return failStep("Assembling Recordsets List", "Unable to display recordsets list")
		}
		results = string(rjson)
	} else {
		results = renderRecordsetListTable(resp.RecordSets)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordsets List ... %s\n", color.GreenString("[OK]"))

	// Write output to file or console
	if len(outputPath) > 1 {
		rlfHandle, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "Failed to create output file. Error: %s", err.Error())
		}
		defer func() { _ = rlfHandle.Close() }()
		_, err = rlfHandle.WriteString(string(results))
		if err != nil {
			return failStep("Writing Output", "Unable to write zone list output to file")
		}
		_ = rlfHandle.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	} else {
		_, _ = fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
