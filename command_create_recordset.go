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

func cmdCreateRecordset(c *cli.Context) error {
	//Validate zone name argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	//Initialize Edgegrid session and DNS client
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	fmt.Println("All args:", os.Args)
	fmt.Println("Parsed flag names:", c.FlagNames())
	fmt.Println("--rdata slice:", c.StringSlice("rdata"))
	fmt.Println("IsSet(\"rdata\"):", c.IsSet("rdata"))
	fmt.Println("c.String(\"rdata\"):", c.String("rdata"))

	var (
		zonename   string
		outputPath string
		inputPath  string
		// json
		// suppress
	)

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
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	fmt.Println("Preparing recordset ", "")

	newrecord := &dns.RecordBody{}

	rdata := c.StringSlice("rdata")
	fmt.Println("rdata values:", rdata)

	missingFlags := []string{}

	if !c.IsSet("name") {
		missingFlags = append(missingFlags, "--name")
	}
	if !c.IsSet("type") {
		missingFlags = append(missingFlags, "--type")
	}
	if !c.IsSet("ttl") {
		missingFlags = append(missingFlags, "--ttl")
	}
	if !c.IsSet("rdata") {
		missingFlags = append(missingFlags, "--rdata")
	}

	if len(missingFlags) > 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(
			color.RedString("Missing required flags: %s", strings.Join(missingFlags, ", ")),
			1,
		)
	}

	//Load recordset from JSOn file if provided
	if c.IsSet("file") {
		data, err := os.ReadFile(filepath.FromSlash(inputPath))
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to read input file"), 1)
		}
		recordset := &dns.RecordSet{}
		err = json.Unmarshal(data, recordset)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to parse json file content into recordset"), 1)
		}
		newrecord.Name = recordset.Name
		newrecord.RecordType = recordset.Type
		newrecord.TTL = recordset.TTL
		newrecord.Target = recordset.Rdata

		// Construct recordset from CLI flags
	} else if c.IsSet("type") {
		if !c.IsSet("name") || !c.IsSet("ttl") || !c.IsSet("rdata") {
			cli.ShowCommandHelp(c, c.Command.Name)
			return cli.NewExitError(color.RedString("Field flags missing for recordset creation"), 1)
		}
		newrecord.RecordType = strings.ToUpper(c.String("type"))
		newrecord.Name = c.String("name")
		newrecord.TTL = c.Int("ttl")
		newrecord.Target = c.StringSlice("rdata")
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset field values or input file are required"), 1)
	}

	fmt.Println("Flag values passed:")
	fmt.Printf("  name: %s\n", c.String("name"))
	fmt.Printf("  type: %s\n", c.String("type"))
	fmt.Printf("  ttl: %d\n", c.Int("ttl"))
	fmt.Printf("  rdata: %v\n", c.StringSlice("rdata"))

	// Check if record already exists
	existing, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: newrecord.RecordType,
		Name:       newrecord.Name,
	})
	if err == nil && existing.RecordType != "" {
		return cli.NewExitError(color.RedString("Recordset already exists"), 1)
	} /*else {
		if !dns.ConfigDNSError() || !err.(dns.ConfigDNSError).NotFound() {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failure while checking recordset existance. Error: %s", err.Error())), 1)
		}
	}*/

	// Create new recordset
	fmt.Println("Creating Recordset  ", "")
	err = dnsClient.CreateRecord(ctx, dns.CreateRecordRequest{Zone: zonename, Record: newrecord})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset create failed. Error: %s", err)), 1)
	}

	// Retrieve recordset after creation
	fmt.Println("Verifying Recordset  ", "")
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{Zone: zonename, RecordType: newrecord.RecordType, Name: newrecord.Name})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read recordset content. Error: %s", err.Error())), 1)
	}

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	results := ""
	fmt.Println(color.BlueString("Assembling recordset Content ", ""))

	// Format output as JSON or table
	if c.IsSet("json") && c.Bool("json") {
		recordset := &dns.RecordSet{}
		recordset.Name = record.Name
		recordset.Type = record.RecordType
		recordset.TTL = record.TTL
		recordset.Rdata = record.Target
		zjson, err := json.MarshalIndent(recordset, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to marshal recordset"), 1)
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, record)
	}

	// Write to file if output path is specified
	if len(outputPath) > 1 {
		fmt.Printf("Writing Output to %s ", outputPath)

		rsHandle, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer rsHandle.Close()
		_, err = rsHandle.WriteString(string(results))
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to write zone output to file"), 1)
		}
		rsHandle.Sync()
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
