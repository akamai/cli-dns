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
	"strings"

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveZone(c *cli.Context) error {

	// Initialize context and EdgeGrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Validate zonename argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	zonename := c.Args().First()

	fmt.Fprintf(c.App.Writer, "Fetching zone...")

	// Fetch zone details
	zoneResp, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{
		Zone: zonename,
	})
	if err != nil {
		return cli.NewExitError(color.RedString("Zone not found "), 1)
	}

	fmt.Fprintln(c.App.Writer, fmt.Sprintf(" [%s]", color.GreenString("OK")))

	if strings.EqualFold(zoneResp.Type, "ALIAS") {
		// Print zone details only
		fmt.Fprintln(c.App.Writer, "")
		renderZoneTable(zoneResp, nil, c)
		return nil
	}

	// Fetch all recordsets for the zone
	recordsResp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{Zone: zonename})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("failed to retrieve zone: %s", err)), 1)
	}

	filterSlice := c.StringSlice("filter")
	filter := make(map[string]bool)
	for _, recordType := range filterSlice {
		filter[strings.ToUpper(recordType)] = true
	}

	// Filter recordsets based on record type
	filteredRecords := []dns.RecordSet{}
	for _, rec := range recordsResp.RecordSets {
		if len(filter) == 0 || filter[strings.ToUpper(rec.Type)] {
			filteredRecords = append(filteredRecords, rec)
		}
	}

	// Output as JSON
	if c.IsSet("json") && c.Bool("json") {
		jsonObj := map[string]interface{}{
			"zone":    zoneResp,
			"records": filteredRecords,
		}
		out, err := json.MarshalIndent(jsonObj, "", " ")
		if err != nil {
			return cli.NewExitError("failed to marshal JSON output", 1)
		}
		fmt.Fprintln(c.App.Writer, string(out))
		return nil
	}

	fmt.Fprintln(c.App.Writer, "")
	renderZoneTable(zoneResp, filteredRecords, c)

	return nil

}
