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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveRecordset(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename   string
		outputPath string
		name       string
		rstype     string
		// json
		// suppress
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

	rstype = c.String("type")
	name = c.String("name")

	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	// Single recordset ops use RecordBody as return Object
	akamai.StartSpinner("Retrieving Recordset  ", "")
	record, err := dnsv2.GetRecord(zonename, name, rstype)
	if err != nil {
		if dnsv2.IsConfigDNSError(err) && err.(dnsv2.ConfigDNSError).NotFound() {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Recordset not found"), 1)
		} else {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read recordset content. Error: %s", err.Error())), 1)
		}
	}

	results := ""
	akamai.StartSpinner("Assembling recordset Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		// output as recordset
		recordset := &dnsv2.Recordset{}
		recordset.Name = record.Name
		recordset.Type = record.RecordType
		recordset.TTL = record.TTL
		recordset.Rdata = record.Target
		zjson, err := json.MarshalIndent(recordset, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to marshal recordset"), 1)
		}
		results = string(zjson)
	} else {
		results = renderRecordsetTable(zonename, record, c)
	}
	akamai.StopSpinnerOk()

	if len(outputPath) > 1 {
		akamai.StartSpinner(fmt.Sprintf("Writing Output to %s ", outputPath), "")
		// pathname and exists?
		rsHandle, err := os.Create(outputPath)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer rsHandle.Close()
		_, err = rsHandle.WriteString(string(results))
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to write zone output to file"), 1)
		}
		rsHandle.Sync()
		akamai.StopSpinnerOk()
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
