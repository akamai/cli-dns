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
	"sort"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateRecordset(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Initialize context and Edgegrid session
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

	// Validate zonename argument
	if c.NArg() == 0 {
		return failStep("Preparing recordset", "zonename is required")
	}

	zonename = c.Args().First()

	// Check if the zone is an ALIAS zone
	zoneResp, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{
		Zone: zonename,
	})
	if err != nil {
		return failStep("Preparing recordset", "Failed to retrieve zone information for %s. Error: %s", zonename, err)
	}
	if strings.EqualFold(zoneResp.Type, "ALIAS") {
		return failStep("Preparing recordset", "Zone %s is an ALIAS zone and does not have recordsets", zonename)
	}

	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	newrecord := &dns.RecordBody{}
	setchange := false
	if c.IsSet("file") {
		// Read and parse the file into a dns.RecordSet
		newrecordset := &dns.RecordSet{}
		data, err := os.ReadFile(filepath.FromSlash(inputPath))
		if err != nil {
			return failStep("Preparing recordset", "Failed to read input file")
		}

		err = json.Unmarshal(data, &newrecordset)
		if err != nil {
			return failStep("Preparing recordset", "Failed to parse json file content into recordset")
		}
		newrecord.Name = newrecordset.Name
		newrecord.RecordType = newrecordset.Type
		if !intPtrEqual(newrecord.TTL, &newrecordset.TTL) {
			setchange = true
		}
		newrecord.TTL = &newrecordset.TTL
		sort.Strings(newrecord.Target)
		sort.Strings(newrecordset.Rdata)
		if !setchange && strings.Join(newrecord.Target, " ") != strings.Join(newrecordset.Rdata, " ") {
			setchange = true
		}
		newrecord.Target = newrecordset.Rdata
	} else if c.IsSet("type") && c.IsSet("name") {
		// Update recordset CLI flags

		newrecord.RecordType = strings.ToUpper(c.String("type"))
		newrecord.Name = c.String("name")
		if c.IsSet("ttl") {
			ttlValue := c.Int("ttl")
			newrecord.TTL = &ttlValue
			setchange = true
		}
		if c.IsSet("rdata") {
			newrecord.Target = c.StringSlice("rdata")
			setchange = true
		}
		fmt.Printf("Preparing recordset ... %s\n", color.GreenString("[OK]"))
	} else {
		return failStep("Preparing recordset", "Recordset field values or input file are required")
	}

	// Retrieve recordset for update
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       newrecord.Name,
		RecordType: newrecord.RecordType,
	})
	if err != nil {
		return failStep("Preparing recordset", "Failure retrieving recordset. Error: %s", err)
	}

	if !c.IsSet("file") {
		if !c.IsSet("ttl") {
			newrecord.TTL = &record.TTL
		}
		if !c.IsSet("rdata") {
			newrecord.Target = record.Target
		}
	}

	if !setchange {
		_, _ = fmt.Fprintln(c.App.Writer, "No recordset change detected")
		return nil
	}

	fmt.Printf("Updating Recordset ... %s\n", color.GreenString("[OK]"))

	// Update recordset
	err = dnsClient.UpdateRecord(ctx, dns.UpdateRecordRequest{
		Zone:   zonename,
		Record: newrecord,
	})
	if err != nil {
		return failStep("Updating Recordset", "Recordset update failed. Error: %s", err.Error())
	}
	fmt.Printf("Updating Recordset ... %s\n", color.GreenString("[OK]"))

	// Fetch updated recordset
	updatedRecord, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       newrecord.Name,
		RecordType: newrecord.RecordType,
	})
	if err != nil {
		return failStep("Verifying Recordset", "Failed to read recordset content. Error: %s", err)
	}
	fmt.Printf("Verifying Recordset ... %s\n", color.GreenString("[OK]"))

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	results := ""
	// Output either as JSON or table format
	if c.IsSet("json") && c.Bool("json") {
		recordset := &dns.RecordSet{}
		recordset.Name = record.Name
		recordset.Type = updatedRecord.RecordType
		recordset.TTL = updatedRecord.TTL
		recordset.Rdata = updatedRecord.Target
		zjson, err := json.MarshalIndent(recordset, "", "  ")
		if err != nil {
			return failStep("Assembling Recordset Content", "Unable to marshal recordset")
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, updatedRecord)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordset Content ... %s\n", color.GreenString("[OK]"))

	// Write output to file or console
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
