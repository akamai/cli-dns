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
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func cmdCreateRecordset(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename   string
		outputPath string
		inputPath  string
		// json
		// suppress
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}
	zonename = c.Args().First()
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	akamai.StartSpinner("Preparing recordset ", "")
	// Single recordset ops use RecordBody as return Object
	newrecord := &dnsv2.RecordBody{}
	if c.IsSet("file") {
		newrecordset := &dnsv2.Recordset{}
		// Read in json file
		data, err := ioutil.ReadFile(inputPath)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to read input file"), 1)
		}
		// set local variables and Object
		err = json.Unmarshal(data, &newrecordset)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to parse json file content into recordset"), 1)
		}
		newrecord.Name = newrecordset.Name
		newrecord.RecordType = newrecordset.Type
		newrecord.TTL = newrecordset.TTL
		newrecord.Target = newrecordset.Rdata
	} else if c.IsSet("type") {
		if !c.IsSet("name") || !c.IsSet("ttl") || !c.IsSet("rdata") {
			akamai.StopSpinnerFail()
			cli.ShowCommandHelp(c, c.Command.Name)
			return cli.NewExitError(color.RedString("Field flags missing for recordset creation"), 1)
		}
		newrecord.RecordType = strings.ToUpper(c.String("type"))
		newrecord.Name = c.String("name")
		newrecord.TTL = c.Int("ttl")
		newrecord.Target = c.StringSlice("rdata")
	} else {
		akamai.StopSpinnerFail()
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Recordset field values or input file are required"), 1)
	}
	// See if already exists
	record, err := dnsv2.GetRecord(zonename, newrecord.Name, newrecord.RecordType) // returns RecordBody!
	if err == nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Recordset already exists"), 1)
	} else {
		if !dnsv2.IsConfigDNSError(err) || !err.(dnsv2.ConfigDNSError).NotFound() {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failure while checking recordset existance. Error: %s", err.Error())), 1)
		}
	}
	akamai.StopSpinnerOk()
	akamai.StartSpinner("Creating Recordset  ", "")
	err = newrecord.Save(zonename, true)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset create failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()
	akamai.StartSpinner("Verifying Recordset  ", "")
	record, err = dnsv2.GetRecord(zonename, newrecord.Name, newrecord.RecordType)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read recordset content. Error: %s", err.Error())), 1)
	}
	// suppress result output?
	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
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

func renderRecordsetTable(zone string, set *dnsv2.RecordBody, c *cli.Context) string {

	outString := "Zone Recordset"
	outString += ""
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetHeader([]string{"NAME", "TYPE", "TTL", "RDATA"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)
	table.SetCaption(true, fmt.Sprintf("Zone: %s", zone))

	if set == nil {
		return outString
	} else {
		name := set.Name
		rstype := set.RecordType
		ttl := strconv.Itoa(set.TTL)
		for i, targ := range set.Target {
			if i == 0 {
				table.Append([]string{name, rstype, ttl, targ})
			} else {
				table.Append([]string{" ", " ", " ", targ})
			}
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
