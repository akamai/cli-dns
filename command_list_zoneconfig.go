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
	"github.com/fatih/color"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/urfave/cli"
)

type ZoneSummary struct {
	Zone            string
	Type            string
	ActivationState string
	ContractId      string
}

type ZoneSummaryList struct {
	Zones []*ZoneSummary
}

type ZoneList struct {
	Zones []*dns.ZoneResponse
}

func cmdListZoneconfig(c *cli.Context) error {

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Build zone list query from CLI flags
	query := dns.ListZonesRequest{
		ShowAll: true,
		SortBy:  "zone",
	}

	if c.IsSet("contractid") {
		contractid := c.String("contractid")
		query.ContractIDs = contractid
	}
	if c.IsSet("type") {
		types := c.StringSlice("type")
		query.Types = strings.Join(types, ",")
	}
	if c.IsSet("search") {
		query.Search = c.String("search")
	}

	// Fetch zones from DNS client
	resp, err := dnsClient.ListZones(ctx, query)
	if err != nil {
		return fmt.Errorf("zone list retrieval failed: %v", err)
	}
	zones := resp.Zones
	var output string

	// Format output in summary or table, optionally as JSON
	if c.Bool("summary") {
		if c.Bool("json") {
			summaryList := ZoneSummaryList{}
			for _, z := range zones {
				summaryList.Zones = append(summaryList.Zones, &ZoneSummary{
					Zone:            z.Zone,
					Type:            z.Type,
					ActivationState: z.ActivationState,
					ContractId:      z.ContractID,
				})
			}
			b, err := json.MarshalIndent(summaryList, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal summary JSON: %v", err)
			}
			output = string(b)
		} else {
			output = renderZoneSummaryListTable(zones)
		}
	} else {
		if c.Bool("json") {
			b, err := json.MarshalIndent(zones, "", " ")
			if err != nil {
				return fmt.Errorf("failed to marshal full JSON: %v", err)
			}
			output = string(b)
		} else {
			output = renderZoneListTable(zones)
		}
	}

	// Write output to file or print to console
	if outFile := c.String("output"); outFile != "" {
		path := filepath.FromSlash(outFile)
		if err := os.WriteFile(path, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write to output file %v", err)
		}
		fmt.Fprintln(c.App.Writer, color.GreenString("Output written to %s", path))
	} else {
		fmt.Fprintln(c.App.Writer, output)
	}

	return nil

}
