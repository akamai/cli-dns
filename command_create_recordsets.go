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
	"strings"

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCreateRecordsets(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Initialize context and EdgeGrid session
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

	// Validate zone name argument
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
		return failStep("Preparing recordsets", "Zone %s is an ALIAS zone and cannot have recordsets", zonename)
	}

	// Get input and output file paths if set
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
	fmt.Printf("Preparing recordsets ... %s\n", color.GreenString("[OK]"))

	// Read and parse input JSON file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return failStep("Fetching Recordset data", "Failed to read input file")
	}
	fmt.Printf("Fetching Recordset data ... %s\n", color.GreenString("[OK]"))

	var wrapper struct {
		RecordSets []dns.RecordSet `json:"recordsets"`
	}
	err = json.Unmarshal(data, &wrapper)
	if err != nil {
		return failStep("Fetching Recordset data", "Failed to parse json file content: %s", err)
	}

	// Create multiple recordsets
	req := dns.CreateRecordSetsRequest{
		Zone: zonename,
		RecordSets: &dns.RecordSets{
			RecordSets: wrapper.RecordSets,
		},
	}

	if err := dnsClient.CreateRecordSets(ctx, req); err != nil {
		return failStep("Creating Recordsets", "Failed to create recordset: %v", err)
	}

	fmt.Printf("Creating Recordsets ... %s\n", color.GreenString("[OK]"))

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	// Retrieve updated list of recordsets
	fmt.Fprintf(os.Stderr, "Retrieving Recordsets List ... %s\n", color.GreenString("[OK]"))
	resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{Zone: zonename, QueryArgs: &dns.RecordSetQueryArgs{ShowAll: true}})
	if err != nil {
		return failStep("Retrieving Recordsets List", "Recordset List retrieval failed. Error: %s", err)
	}

	// Format recordsets for output
	recordsetList := RecordsetList{Recordsets: resp.RecordSets}
	results := ""
	if c.IsSet("json") && c.Bool("json") {
		rjson, err := json.MarshalIndent(recordsetList, "", "  ")
		if err != nil {
			return failStep("Assembling Recordsets List", "Unable to display recordsets list")
		}
		results = string(rjson)
	} else {
		results = renderRecordsetListTable(recordsetList.Recordsets)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordsets List ... %s\n", color.GreenString("[OK]"))

	// Write to file if output path is specified or print to stdout
	if len(outputPath) > 1 {
		//fmt.Println(color.GreenString("Writing Output to %s", outputPath))
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
		fmt.Println(color.GreenString("Output written to %s", outputPath))
		return nil
	} else {
		_, _ = fmt.Fprintln(c.App.Writer, results)
	}

	return nil

}
