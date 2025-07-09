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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdAddRecord(c *cli.Context) error {

	//Validate postional arguments; record type and zone name
	if c.NArg() < 2 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("record type and zonename are required"), 1)
	}

	recordType := strings.ToUpper(c.Args().Get(0))
	zonename := strings.TrimSuffix(c.Args().Get(1), ".")

	//validate required flags
	if !c.IsSet("name") || !c.IsSet("rdata") || !c.IsSet("ttl") {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("--name, --rdata and --ttl are required"), 1)
	}

	name := c.String("name")
	if !strings.HasSuffix(name, "."+zonename) {
		name = fmt.Sprintf("%s.%s", name, zonename)
	}
	if !strings.HasSuffix(name, "."+zonename) {
		return cli.NewExitError(color.RedString("record name must be within the zone %s", zonename), 1)
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
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

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

	// Define new record
	newrecord := &dns.RecordBody{
		RecordType: recordType,
		Name:       name,
		TTL:        ttl,
		Target:     rdata,
	}

	// Check if record already exists
	existing, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: newrecord.RecordType,
		Name:       newrecord.Name,
	})

	if err == nil && existing.RecordType != "" {

		fmt.Println("Record already exists, updating it instead...")

		//Merge TTL and RDATA values if needed
		ttlChanged := existing.TTL != newrecord.TTL

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
			fmt.Println(color.BlueString("No changes to update."))
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
			return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset update failed. Error: %s", err)), 1)
		}
	} else {
		// Create a new record
		fmt.Println(color.BlueString("Creating new recordset..."))
		err = dnsClient.CreateRecord(ctx, dns.CreateRecordRequest{
			Zone:   zonename,
			Record: newrecord,
		})
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset create failed. Error: %s", err)), 1)
		}
	}

	// Retrieve record after creation/update
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: newrecord.RecordType,
		Name:       newrecord.Name,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read recordset content. Error: %s", err.Error())), 1)
	}
	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	//Output the recordset
	fmt.Println(color.BlueString("Assembling recordset Content... ", ""))
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
			return cli.NewExitError(color.RedString("Unable to marshal recordset"), 1)
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, record)
	}

	if len(outputPath) > 1 {
		//fmt.Println(color.GreenString("Writing output to %s", outputPath))
		rsHandle, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer rsHandle.Close()
		_, err = rsHandle.WriteString(results)
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to write zone output to file"), 1)
		}
		rsHandle.Sync()
		fmt.Println(color.GreenString("Output written to %s", outputPath))
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
