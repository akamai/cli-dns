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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdDeleteRecordset(c *cli.Context) error {

	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}
	zonename := c.Args().First()

	if !c.IsSet("name") || !c.IsSet("type") {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset name and type field values are required"), 1)
	}
	recordType := c.String("type")
	recordName := c.String("name")

	fmt.Println("Checking Recordset existance  ", "")
	// See if already exists
	_, err = dnsClient.GetRecord(ctx, dns.GetRecordRequest{
		Zone:       zonename,
		Name:       recordName,
		RecordType: recordType,
	}) // returns RecordBody!
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failure retrieving recordset. Error: %s", err)), 1)
	}

	// Single recordset ops use RecordBody as return Object
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
