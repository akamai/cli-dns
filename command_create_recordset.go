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

func cmdCreateRecordset(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Validate zone name argument
	if c.NArg() == 0 {
		return failStep("Preparing recordset", "zonename is required")
	}

	// Initialize Edgegrid session and DNS client
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing recordset", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	var (
		zonename   string
		outputPath string
		inputPath  string
	)

	zonename = c.Args().First()
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

	// Get input and output file paths if set
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	newrecord := &dns.RecordBody{}

	// Load recordset from JSON file if provided
	if c.IsSet("file") {
		data, err := os.ReadFile(filepath.FromSlash(inputPath))
		if err != nil {
			return failStep("Preparing recordset", "Failed to read input file")
		}
		recordset := &dns.RecordSet{}
		err = json.Unmarshal(data, recordset)
		if err != nil {
			return failStep("Preparing recordset", "Failed to parse json file content into recordset")
		}
		newrecord.Name = recordset.Name
		newrecord.RecordType = recordset.Type
		newrecord.TTL = &recordset.TTL
		newrecord.Target = recordset.Rdata
	} else if c.IsSet("type") {
		if !c.IsSet("name") || !c.IsSet("ttl") || !c.IsSet("rdata") {
			return failStep("Preparing recordset", "Field flags missing for recordset creation")
		}
		newrecord.RecordType = strings.ToUpper(c.String("type"))
		newrecord.Name = c.String("name")
		ttlValue := c.Int("ttl")  // 1. Save to a named variable first
		newrecord.TTL = &ttlValue // 2. Now you can safely take the address
		newrecord.Target = c.StringSlice("rdata")
	} else {
		return failStep("Preparing recordset", "Recordset field values or input file are required")
	}
	fmt.Printf("Preparing recordset ... %s\n", color.GreenString("[OK]"))

	// Check if record already exists
	existing, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: newrecord.RecordType,
		Name:       newrecord.Name,
	})
	if err == nil && existing.RecordType != "" {
		return failStep("Checking Recordset Existence", "Recordset already exists")
	}
	fmt.Printf("Checking Recordset Existence ... %s\n", color.GreenString("[OK]"))

	// Create new recordset
	err = dnsClient.CreateRecord(ctx, dns.CreateRecordRequest{Zone: zonename, Record: newrecord})
	if err != nil {
		return failStep("Creating Recordset", "Recordset create failed. Error: %s", err)
	}
	fmt.Printf("Creating Recordset ... %s\n", color.GreenString("[OK]"))

	// Retrieve recordset after creation
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{Zone: zonename, RecordType: newrecord.RecordType, Name: newrecord.Name})
	if err != nil {
		return failStep("Verifying Recordset", "Failed to read recordset content. Error: %s", err.Error())
	}
	fmt.Printf("Verifying Recordset ... %s\n", color.GreenString("[OK]"))

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	results := ""
	// Format output as JSON or table
	if c.IsSet("json") && c.Bool("json") {
		recordset := &dns.RecordSet{}
		recordset.Name = record.Name
		recordset.Type = record.RecordType
		recordset.TTL = record.TTL
		recordset.Rdata = record.Target
		zjson, err := json.MarshalIndent(recordset, "", "  ")
		if err != nil {
			return failStep("Assembling Recordset Content", "Unable to marshal recordset")
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, record)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordset Content ... %s\n", color.GreenString("[OK]"))

	// Write to file if output path is specified
	if len(outputPath) > 1 {
		rsHandle, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "Failed to create output file. Error: %s", err.Error())
		}
		defer func() { _ = rsHandle.Close() }()
		_, err = rsHandle.WriteString(string(results))
		if err != nil {
			return failStep("Writing Output", "Unable to write output to file")
		}
		_ = rsHandle.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	}

	_, _ = fmt.Fprintln(c.App.Writer, results)
	return nil
}
