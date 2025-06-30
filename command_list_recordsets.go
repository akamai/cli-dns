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

type RecordsetList struct {
	Recordsets []dns.RecordSet
}

func cmdListRecordsets(c *cli.Context) error {

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename required"), 1)
	}

	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	zonename := c.Args().First()

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

	fmt.Fprintln(os.Stderr, color.BlueString("Retrieving Recordsets List"))
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

	resp, err := dnsClient.GetRecordSets(ctx, req)
	if err != nil {
		return cli.NewExitError(color.RedString("Recordset List retrieval failed %s", err), 1)
	}

	recordsets := resp.RecordSets
	var results string

	if c.Bool("json") {
		output := RecordsetList{Recordsets: recordsets}
		b, err := json.MarshalIndent(output, "", " ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to format JSON"), 1)
		}
		results = string(b)
	} else {
		results = renderRecordsetListTable(zonename, recordsets)
	}

	if outputPath != "" {
		f, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to create output file: %s", err), 1)
		}

		defer f.Close()
		f.WriteString(results)
		f.Sync()
		fmt.Fprintf(os.Stderr, "Output is written to %s", outputPath)
	}

	fmt.Fprintln(c.App.Writer, "")
	fmt.Fprintln(c.App.Writer, results)
	return nil
}
