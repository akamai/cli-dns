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
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRmRecord(c *cli.Context) error {
	// Initialize context and EdgeGrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing recordset", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Validate record type and zone name arguments
	if c.NArg() < 2 {
		return failStep("Preparing recordset", "record type and zonename are required")
	}
	recordType := strings.ToUpper(c.Args().Get(0))
	zonename := c.Args().Get(1)

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

	if !c.IsSet("name") {
		return failStep("Preparing recordset", "Record name (--name) is required")
	}
	fmt.Printf("Preparing recordset ... %s\n", color.GreenString("[OK]"))
	name := c.String("name")

	fqdn := name
	if !strings.HasSuffix(name, "."+zonename) {
		fqdn = name + "." + zonename
	}

	// Get list of recordsets matching the type and zone
	listResp, err := dnsClient.GetRecordList(ctx, dns.GetRecordListRequest{
		Zone:       zonename,
		RecordType: recordType,
	})
	if err != nil {
		return failStep("Looking up records to delete", "Failed to list records: %s", err)
	}

	// Filter matching records by name
	matching := []dns.RecordBody{}
	for _, rec := range listResp.RecordSets {
		if rec.Name == fqdn {
			matching = append(matching, dns.RecordBody{
				Name:       rec.Name,
				RecordType: rec.Type,
				TTL:        &rec.TTL,
				Target:     rec.Rdata,
			})
		}
	}

	if len(matching) == 0 {
		return failStep("Looking up records to delete", "No matching records found.")
	}

	// If multiple records match, ask user unless --force-multiple is set
	if len(matching) > 1 && !c.Bool("force-multiple") {
		if c.Bool("non-interactive") {
			return failStep("Looking up records to delete", "Multiple records found. Use --force-multiple in non-interactive mode.")
		}
		fmt.Printf("Looking up records to delete ... %s\n", color.GreenString("[OK]"))
		fmt.Printf("Multiple records matched for %s %s:\n", recordType, fqdn)
		for _, rec := range matching {
			fmt.Printf("- TTL: %d, RDATA: %v\n", intValue(rec.TTL), rec.Target)
		}
		fmt.Print("Are you sure you want to delete all matching records? [y/N]: ")
		reader := bufio.NewReader(os.Stdin)
		resp, _ := reader.ReadString('\n')
		resp = strings.ToLower(strings.TrimSpace(resp))
		if resp != "y" && resp != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}
	fmt.Printf("Looking up records to delete ... %s\n", color.GreenString("[OK]"))
	// Delete each matching record
	for _, rec := range matching {
		err = dnsClient.DeleteRecord(ctx, dns.DeleteRecordRequest{
			Zone:       zonename,
			Name:       rec.Name,
			RecordType: rec.RecordType,
		})
		if err != nil {
			return failStep("Deleting Recordset", "Failed to delete record %s %s: %s", rec.RecordType, rec.Name, err)
		}
	}
	fmt.Printf("Deleting Recordset ... %s\n", color.GreenString("[OK]"))
	return nil
}
