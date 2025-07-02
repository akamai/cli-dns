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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdDeleteRecordset(c *cli.Context) error {

	// Initialize context and Edgegrid session
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

	// Validate required flags
	if !c.IsSet("name") || !c.IsSet("type") {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset name and type field values are required"), 1)
	}
	recordType := c.String("type")
	recordName := c.String("name")

	// Check if recordset exists
	fmt.Println("Checking Recordset existance  ", "")

	_, err = dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       recordName,
		RecordType: recordType,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failure retrieving recordset. Error: %s", err)), 1)
	}

	// Delete recordset
	err = dnsClient.DeleteRecord(ctx, dns.DeleteRecordRequest{
		Zone:       zonename,
		Name:       recordName,
		RecordType: recordType,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("failed to delete record: %s", err)), 1)
	}

	fmt.Println(color.GreenString("Record Deleted Successfully"))
	return nil
}
