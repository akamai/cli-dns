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
	"github.com/olekukonko/tablewriter"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveZone(c *cli.Context) error {

	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("hostname is required"), 1)
	}

	zonename := c.Args().First()

	fmt.Fprintf(c.App.Writer, "Fetching zone...")

	zoneResp, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{
		Zone: zonename,
	})
	if err != nil {
		return cli.NewExitError(color.RedString("Zone not found "), 1)
	}

	fmt.Fprintln(c.App.Writer, fmt.Sprintf(" [%s]", color.GreenString("OK")))

	recordsResp, err := dnsClient.GetRecordSets(ctx, dns.GetRecordSetsRequest{Zone: zonename})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("failed to retrieve zone: %s", err)), 1)
	}

	filterSlice := c.StringSlice("filter")
	filter := make(map[string]bool)
	for _, recordType := range filterSlice {
		filter[strings.ToUpper(recordType)] = true
	}

	filteredRecords := []dns.RecordSet{}
	for _, rec := range recordsResp.RecordSets {
		if len(filter) == 0 || filter[strings.ToUpper(rec.Type)] {
			filteredRecords = append(filteredRecords, rec)
		}
	}

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

func renderZoneTable(zone *dns.GetZoneResponse, records []dns.RecordSet, c *cli.Context) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)

	table.SetHeader([]string{"Field", "value"})
	table.Append([]string{"Zone", zone.Zone})
	table.Append([]string{"Type", zone.Type})
	table.Append([]string{"Masters", strings.Join(zone.Masters, ", ")})
	table.Append([]string{"Comment", zone.Comment})
	table.Append([]string{"Contract ID", zone.ContractID})
	table.Append([]string{"SignAndServe", fmt.Sprintf("%v", zone.SignAndServeAlgorithm)})
	table.Append([]string{"Target", zone.Target})
	table.Append([]string{"EndCustomerID", zone.EndCustomerID})
	table.Append([]string{"Activation State", zone.ActivationState})
	table.Append([]string{"Last Modified By", zone.LastModifiedBy})
	table.Append([]string{"Last Modified Date", zone.LastModifiedDate})
	table.Append([]string{"Version ID", zone.VersionID})

	if zone.TSIGKey != nil {
		table.Append([]string{"TSIG Name", zone.TSIGKey.Name})
		table.Append([]string{"TSIG Algorithm", zone.TSIGKey.Algorithm})
		table.Append([]string{"TSIG Secret", zone.TSIGKey.Secret})
	}

	table.Render()
	fmt.Fprintln(c.App.Writer, tableString.String())

	if len(records) > 0 {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, "DNS Records: ")
		fmt.Fprintln(c.App.Writer, "")

		recordsTableString := &strings.Builder{}
		recordsTable := tablewriter.NewWriter(recordsTableString)
		recordsTable.SetHeader([]string{"Name", "Type", "TTL", "Data"})
		recordsTable.SetAutoWrapText(false)
		recordsTable.SetRowLine(true)

		for _, rec := range records {
			for _, data := range rec.Rdata {
				recordsTable.Append([]string{
					rec.Name, rec.Type, fmt.Sprintf("%d", rec.TTL), data,
				})
			}
		}
		recordsTable.Render()
		fmt.Fprintln(c.App.Writer, recordsTableString.String())
	}
}
