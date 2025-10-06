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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCreateRecordsets(c *cli.Context) error {

	// Initialize context and EdgeGrid session
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

	// Validate zone name argument
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
		return cli.NewExitError(color.RedString(fmt.Sprintf("Zone %s is an ALIAS zone and cannot have recordsets", zonename)), 1)
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
		return cli.NewExitError(color.RedString("Input file is required"), 1)
	}
	fmt.Println("Fetching Recordset data ", "")

	// Read and parse input JSON file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to read input file"), 1)
	}

	var wrapper struct {
		RecordSets []dns.RecordSet `json:"recordsets"`
	}
	err = json.Unmarshal(data, &wrapper)
	if err != nil {
		return cli.NewExitError(color.RedString("Failed to parse json file content: %s", err), 1)
	}

	// Create multiple recordsets
	req := dns.CreateRecordSetsRequest{
		Zone: zonename,
		RecordSets: &dns.RecordSets{
			RecordSets: wrapper.RecordSets,
		},
	}

	if err := dnsClient.CreateRecordSets(ctx, req); err != nil {
		return cli.NewExitError(fmt.Sprintf("Failed to create recordset: %v", err), 1)
	}

	fmt.Println("Creating Recordsets ", "")

	if c.IsSet("suppress") && c.Bool("suppress") {

		return nil

	}

	// Retrieve updated list of recordsets
	fmt.Println(color.BlueString("Retrieving Full Recordsets List... ", ""))
	resp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{Zone: zonename, QueryArgs: &dns.RecordSetQueryArgs{ShowAll: true}})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset List retrieval failed. Error: %s", err)), 1)
	}

	// Format recordsets for output
	recordsetList := RecordsetList{Recordsets: resp.RecordSets}
	results := ""
	fmt.Println(color.BlueString("Assembling Recordsets List... ", ""))
	if c.IsSet("json") && c.Bool("json") {
		rjson, err := json.MarshalIndent(recordsetList, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to display recordsets list"), 1)
		}
		results = string(rjson)
	} else {
		results = renderRecordsetListTable(zonename, recordsetList.Recordsets)
	}

	// Write to file if output path is specified or print to stdout
	if len(outputPath) > 1 {
		//fmt.Println(color.GreenString("Writing Output to %s", outputPath))
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
		fmt.Println(color.GreenString("Output written to %s", outputPath))
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil

}
