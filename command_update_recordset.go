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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateRecordset(c *cli.Context) error {
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

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	zonename = c.Args().First()
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	fmt.Println("Preparing recordset ", "")
	// Single recordset ops use RecordBody as return Object
	newrecord := &dns.RecordBody{}
	setchange := false
	if c.IsSet("file") {
		newrecordset := &dns.RecordSet{}
		// Read in json file
		data, err := os.ReadFile(filepath.FromSlash(inputPath))
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to read input file"), 1)
		}
		// set local variables and Object
		err = json.Unmarshal(data, &newrecordset)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to parse json file content into recordset"), 1)
		}
		newrecord.Name = newrecordset.Name
		newrecord.RecordType = newrecordset.Type
		if newrecord.TTL != newrecordset.TTL {
			setchange = true
		}
		newrecord.TTL = newrecordset.TTL
		sort.Strings(newrecord.Target)
		sort.Strings(newrecordset.Rdata)
		if !setchange && strings.Join(newrecord.Target, " ") != strings.Join(newrecordset.Rdata, " ") {
			setchange = true
		}
		newrecord.Target = newrecordset.Rdata
	} else if c.IsSet("type") && c.IsSet("name") {
		newrecord.RecordType = strings.ToUpper(c.String("type"))
		newrecord.Name = c.String("name")
		if c.IsSet("ttl") {
			newrecord.TTL = c.Int("ttl")
			setchange = true
		}
		if c.IsSet("rdata") {
			newrecord.Target = c.StringSlice("rdata")
			setchange = true
		}
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset field values or input file are required"), 1)
	}

	// See if already exists
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       newrecord.Name,
		RecordType: newrecord.RecordType,
	}) // returns RecordBody!
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failure retrieving recordset. Error: %s", err)), 1)
	}

	// Overlay changed fields
	if !c.IsSet("file") {
		if !c.IsSet("ttl") {
			newrecord.TTL = record.TTL
		}
		if !c.IsSet("rdata") {
			newrecord.Target = record.Target
		}
	}

	if !setchange {
		fmt.Fprintln(c.App.Writer, "No recordset change detected")
		return nil
	}

	fmt.Println("Updating Recordset  ", "")

	err = dnsClient.UpdateRecord(ctx, dns.UpdateRecordRequest{
		Zone:   zonename,
		Record: newrecord,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset update failed. Error: %s", err.Error())), 1)
	}

	fmt.Println("Verifying Recordset  ", "")

	updatedRecord, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       newrecord.Name,
		RecordType: newrecord.RecordType,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read recordset content. Error: %s", err)), 1)
	}
	// suppress result output?
	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	results := ""
	fmt.Println("Assembling recordset Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		// output as recordset
		recordset := &dns.RecordSet{}
		recordset.Name = record.Name
		recordset.Type = updatedRecord.RecordType
		recordset.TTL = updatedRecord.TTL
		recordset.Rdata = updatedRecord.Target
		zjson, err := json.MarshalIndent(recordset, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to marshal recordset"), 1)
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, updatedRecord)
	}

	if len(outputPath) > 1 {
		fmt.Printf("Writing Output to %s ", outputPath)
		// pathname and exists?
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
