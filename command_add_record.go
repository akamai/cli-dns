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
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdAddRecord(c *cli.Context) error {
	//Validate postional arguments; record type and zone name
	if c.NArg() < 2 {
		return failStep("Preparing recordset", "record type and zonename are required")
	}

	recordType := strings.ToUpper(c.Args().Get(0))
	zonename := strings.TrimSuffix(c.Args().Get(1), ".")

	//validate required flags
	if !c.IsSet("name") || !c.IsSet("rdata") || !c.IsSet("ttl") {
		return failStep("Preparing recordset", "--name, --rdata and --ttl are required")
	}

	name := c.String("name")
	if !strings.HasSuffix(name, "."+zonename) {
		name = fmt.Sprintf("%s.%s", name, zonename)
	}
	if !strings.HasSuffix(name, "."+zonename) {
		return failStep("Preparing recordset", "record name must be within the zone %s", zonename)
	}

	ttl := c.Int("ttl")
	rdata := c.StringSlice("rdata")
	outputPath := ""
	if c.IsSet("output") {
		outputPath = filepath.FromSlash(c.String("output"))
	}

	//Set up Edgegrid session and DNS client
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
	fmt.Printf("Preparing recordset ... %s\n", color.GreenString("[OK]"))

	// Define new record
	newrecord := &dns.RecordBody{
		RecordType: recordType,
		Name:       name,
		TTL:        &ttl,
		Target:     rdata,
	}

	// Check if record already exists
	existing, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: newrecord.RecordType,
		Name:       newrecord.Name,
	})

	if err == nil && existing.RecordType != "" {
		//Merge TTL and RDATA values if needed
		ttlChanged := existing.TTL != intValue(newrecord.TTL)

		rdataMap := map[string]bool{}
		for _, r := range existing.Target {
			rdataMap[r] = true
		}
		for _, r := range newrecord.Target {
			rdataMap[r] = true
		}
		mergedRdata := []string{}
		for r := range rdataMap {
			mergedRdata = append(mergedRdata, r)
		}
		sort.Strings(mergedRdata)

		changed := ttlChanged || strings.Join(existing.Target, "") != strings.Join(mergedRdata, "")
		if !changed {
			_, _ = fmt.Fprintln(c.App.Writer, "No recordset change detected")
			return nil
		}

		// Update record with merged RDATA and TTL if record already exists
		updateRecord := &dns.RecordBody{
			Name:       existing.Name,
			RecordType: existing.RecordType,
			TTL:        newrecord.TTL,
			Target:     mergedRdata,
		}

		err = dnsClient.UpdateRecord(ctx, dns.UpdateRecordRequest{
			Zone:   zonename,
			Record: updateRecord,
		})
		if err != nil {
			return failStep("Updating Recordset", "Recordset update failed. Error: %s", err)
		}
		fmt.Printf("Updating Recordset ... %s\n", color.GreenString("[OK]"))
	} else {
		// Create a new record

		err = dnsClient.CreateRecord(ctx, dns.CreateRecordRequest{
			Zone:   zonename,
			Record: newrecord,
		})
		if err != nil {
			return failStep("Creating Recordset", "Recordset create failed. Error: %s", err)
		}
		fmt.Printf("Creating Recordset ... %s\n", color.GreenString("[OK]"))
	}

	// Retrieve record after creation/update
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: newrecord.RecordType,
		Name:       newrecord.Name,
	})
	if err != nil {
		return failStep("Verifying Recordset", "Failed to read recordset content. Error: %s", err.Error())
	}
	fmt.Printf("Verifying Recordset ... %s\n", color.GreenString("[OK]"))
	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	//Output the recordset

	var results string
	if c.IsSet("json") && c.Bool("json") {
		recordset := &dns.RecordSet{
			Name:  record.Name,
			Type:  record.RecordType,
			TTL:   record.TTL,
			Rdata: record.Target,
		}
		zjson, err := json.MarshalIndent(recordset, "", "  ")
		if err != nil {
			return failStep("Assembling Recordset Content", "Unable to marshal recordset")
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, record)
	}
	fmt.Fprintf(os.Stderr, "Assembling Recordset Content ... %s\n", color.GreenString("[OK]"))

	if len(outputPath) > 1 {
		rsHandle, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "Failed to create output file. Error: %s", err.Error())
		}
		defer func() { _ = rsHandle.Close() }()
		_, err = rsHandle.WriteString(results)
		if err != nil {
			return failStep("Writing Output", "Unable to write zone output to file")
		}
		_ = rsHandle.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
	} else {
		_, _ = fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
