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

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveRecordset(c *cli.Context) error {

	// Validate zonename argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}
	zonename := c.Args().First()

	// Validate required flags
	if !c.IsSet("name") || !c.IsSet("type") {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset name and type are required"), 1)
	}

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
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	fmt.Fprintf(os.Stderr, color.BlueString("\n Retrieving Recordset"))
	record, err := dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		RecordType: rstype,
		Name:       name,
	})
	if err != nil {
		if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
			return cli.NewExitError(color.RedString("Recordset not found"), 1)
		}
		return cli.NewExitError(color.RedString("Failed to retrieve recordset: %s", err), 1)
	}

	fmt.Fprintf(os.Stderr, color.BlueString("Assembling Recordset Output...\n"))
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
			return cli.NewExitError(color.RedString("Unable to format JSON: %s", err), 1)
		}
		results = string(b)
	} else {
		results = renderRecordsetTable(zonename, record)
	}
	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to create output file: %s", err), 1)
		}
		defer f.Close()
		f.WriteString(results)
		f.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	}

	fmt.Fprintln(c.App.Writer, "")
	fmt.Fprintln(c.App.Writer, results)
	return nil

}
