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
	"fmt"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdDeleteRecordset(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename string
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}
	zonename = c.Args().First()

	if !c.IsSet("name") || !c.IsSet("type") {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset name and type field values are required"), 1)
	}
	record := dnsv2.RecordBody{}
	record.RecordType = c.String("type")
	record.Name = c.String("name")

	akamai.StartSpinner("Checking Recordset existance  ", "")
	// See if already exists
	_, err = dnsv2.GetRecord(zonename, record.Name, record.RecordType) // returns RecordBody!
	if err != nil {
		if dnsv2.IsConfigDNSError(err) && err.(dnsv2.ConfigDNSError).NotFound() {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Existing recordset not found."), 1)
		} else {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failure retrieving recordset. Error: %s", err.Error())), 1)
		}
	}
	akamai.StopSpinnerOk()
	// Single recordset ops use RecordBody as return Object
	akamai.StartSpinner("Deleting Recordset  ", "")
	err = record.Delete(zonename, true)
	if err != nil {
		if dnsv2.IsConfigDNSError(err) && err.(dnsv2.ConfigDNSError).NotFound() {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Recordset not found"), 1)
		} else {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to delete recordset. Error: %s", err.Error())), 1)
		}
	}
	akamai.StopSpinnerOk()
	return nil
}
