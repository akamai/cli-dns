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

type RecordsetList struct {
	Recordsets []dns.RecordSet
}

func cmdListRecordsets(c *cli.Context) error {
	// Validate zonename argument
	if c.NArg() == 0 {
		return failStep("Preparing recordsets", "zonename required")
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing recordsets", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	zonename := c.Args().First()

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
	fmt.Printf("Preparing recordsets ... %s\n", color.GreenString("[OK]"))

	outputPath := ""
	if c.IsSet("output") {
		outputPath = filepath.FromSlash(c.String("output"))
	}

	typeFilter := c.StringSlice("type")
	search := c.String("search")
	sortby := c.String("sortby")
	if sortby == "" {
		sortby = "type"
	}

	req := dns.GetRecordSetsRequest{
		Zone: zonename,
		QueryArgs: &dns.RecordSetQueryArgs{
			ShowAll: true,
			Search:  search,
			SortBy:  sortby,
		},
	}
	if len(typeFilter) > 0 {
		req.QueryArgs.Types = strings.Join(typeFilter, ",")
	}

	// Fetch recordsets
	resp, err := dnsClient.GetRecordSets(ctx, req)
	if err != nil {
		return failStep("Retrieving Recordsets List", "Recordset List retrieval failed %s", err)
	}
	fmt.Printf("Retrieving Recordsets List ... %s\n", color.GreenString("[OK]"))

	recordsets := resp.RecordSets

	// Format output (JSON or table)
	var results string
	if c.Bool("json") {
		output := RecordsetList{Recordsets: recordsets}
		b, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return failStep("Assembling Recordsets List", "Unable to format JSON")
		}
		results = string(b)
	} else {
		results = renderRecordsetListTable(recordsets)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordsets List ... %s\n", color.GreenString("[OK]"))

	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "Failed to create output file: %s", err)
		}
		defer func() { _ = f.Close() }()
		_, _ = f.WriteString(results)
		if err := f.Sync(); err != nil {
			return failStep("Writing Output", "failed to sync file: %s", err)
		}
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	}

	_, _ = fmt.Fprintln(c.App.Writer, results)
	return nil
}
