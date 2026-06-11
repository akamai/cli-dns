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
	"fmt"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdDeleteRecordset(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing recordset", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Validate zonename argument
	if c.NArg() == 0 {
		return failStep("Preparing recordset", "zonename is required")
	}
	zonename := c.Args().First()

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

	// Validate required flags
	if !c.IsSet("name") || !c.IsSet("type") {
		return failStep("Preparing recordset", "Recordset name and type field values are required")
	}
	recordType := c.String("type")
	recordName := c.String("name")
	fmt.Printf("Preparing recordset ... %s\n", color.GreenString("[OK]"))

	// Check if recordset exists
	_, err = dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       recordName,
		RecordType: recordType,
	})
	if err != nil {
		return failStep("Checking Recordset Existence", "Failure retrieving recordset. Error: %s", err)
	}
	fmt.Printf("Checking Recordset Existence ... %s\n", color.GreenString("[OK]"))

	// Delete recordset
	err = dnsClient.DeleteRecord(ctx, dns.DeleteRecordRequest{
		Zone:       zonename,
		Name:       recordName,
		RecordType: recordType,
	})
	if err != nil {
		return failStep("Deleting Recordset", "failed to delete record: %s", err)
	}
	fmt.Printf("Deleting Recordset ... %s\n", color.GreenString("[OK]"))
	return nil
}
