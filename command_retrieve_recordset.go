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

func cmdRetrieveRecordset(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Validate zonename argument
	if c.NArg() == 0 {
		return failStep("Preparing recordset", "zonename is required")
	}
	zonename := c.Args().First()

	// Validate required flags
	if !c.IsSet("name") || !c.IsSet("type") {
		return failStep("Preparing recordset", "Recordset name and type are required")
	}
	fmt.Printf("Preparing recordset ... %s\n", color.GreenString("[OK]"))

	name := c.String("name")
	rstype := c.String("type")

	outputPath := ""
	if c.IsSet("output") {
		outputPath = filepath.FromSlash(c.String("output"))
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing recordset", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Check if the zone is an ALIAS zone
	zoneResp, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{
		Zone: zonename,
	})
	if err != nil {
		return failStep("Preparing recordset", "Failed to retrieve zone information for %s. Error: %s", zonename, err)
	}
	if strings.EqualFold(zoneResp.Type, "ALIAS") {
		return failStep("Preparing recordset", "Zone %s is an ALIAS zone and cannot have recordsets", zonename)
	}

	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: rstype,
		Name:       name,
	})
	if err != nil {
		if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
			return failStep("Retrieving Recordset", "Recordset not found")
		}
		return failStep("Retrieving Recordset", "Failed to retrieve recordset: %s", err)
	}
	fmt.Printf("Retrieving Recordset ... %s\n", color.GreenString("[OK]"))

	var results string
	if c.Bool("json") {
		rs := &dns.RecordSet{
			Name:  record.Name,
			Type:  record.RecordType,
			TTL:   record.TTL,
			Rdata: record.Target,
		}
		b, err := json.MarshalIndent(rs, "", " ")
		if err != nil {
			return failStep("Assembling Recordset Content", "Unable to format JSON: %s", err)
		}
		results = string(b)
	} else {
		results = renderRecordsetTable(zonename, record)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordset Content ... %s\n", color.GreenString("[OK]"))

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
