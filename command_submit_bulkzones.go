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
	"strings"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func cmdSubmitBulkZones(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		outputPath     string
		contractid     string
		groupid        string
		inputPath      string
		bulkDeleteList *dnsv2.ZoneNameListResponse
		newBulkZones   *dnsv2.BulkZonesCreate
		op             string = "create"
		bypass         bool
	)

	akamai.StartSpinner("Preparing bulk zones submit request ", "")
	queryArgs := dnsv2.ZoneQueryString{}

	if c.IsSet("contractid") {
		contractid = c.String("contractid")
		queryArgs.Contract = contractid
	} else if c.IsSet("create") {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("contractid is required"), 1)
	}
	if c.IsSet("groupid") {
		groupid = c.String("groupid")
		queryArgs.Group = groupid
	}
	if c.IsSet("bypassZoneSafety") && c.Bool("bypassZoneSafety") {
		bypass = true
	}
	if (c.IsSet("create") && c.IsSet("delete")) || (!c.IsSet("create") && !c.IsSet("delete")) {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Either create or delete arg is required. "), 1)
	}
	if c.IsSet("delete") {
		op = "delete"
		bulkDeleteList = &dnsv2.ZoneNameListResponse{}
	} else {
		newBulkZones = &dnsv2.BulkZonesCreate{}
	}
	if op == "create" {
		if bypass {
			fmt.Printf("Warning: bypassZoneSafety arg ignored")
		}
	} else {
		if c.IsSet("contractid") {
			fmt.Printf("Warning: contractid arg ignored")
		}
		if c.IsSet("groupid") {
			fmt.Printf("Warning: groupid arg ignored")
		}
	}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	} else {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(" Bulk create JSON source file must be specified"), 1)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	// Read in json file
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Failed to read input file"), 1)
	}
	// set local variables and Object
	if op == "create" {
		err = json.Unmarshal(data, newBulkZones)
	} else {
		err = json.Unmarshal(data, bulkDeleteList)
	}
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Failed to parse json file content into bulk zones object"), 1)
	}

	akamai.StartSpinner("Submitting Bulk Zones request  ", "")
	//  Submit
	var submitStatus *dnsv2.BulkZonesResponse
	if op == "create" {
		submitStatus, err = dnsv2.CreateBulkZones(newBulkZones, queryArgs)
	} else {
		submitStatus, err = dnsv2.DeleteBulkZones(bulkDeleteList, bypass)
	}
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request submit failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()

	results := ""
	akamai.StartSpinner("Assembling Bulk Zone Response Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(submitStatus, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display zone"), 1)
		}
		results = string(zjson)
	} else {
		results = renderBulkZonesRequestStatusTable(submitStatus, c)
	}
	akamai.StopSpinnerOk()

	if len(outputPath) < 1 {
		// if no output path, write out locally. RequestId has to be retreievable
		outputPath = fmt.Sprintf("bulkSubmitRequest.%s", submitStatus.RequestId)
		outputPath = filepath.FromSlash(outputPath)
	}
	akamai.StartSpinner(fmt.Sprintf("Writing Request Status to %s ", outputPath), "")
	// pathname and exists?
	zfHandle, err := os.Create(outputPath)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
	}
	defer zfHandle.Close()
	outputPath = filepath.FromSlash(outputPath)

	_, err = zfHandle.WriteString(string(results))
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Unable to write zone output to file"), 1)
	}
	zfHandle.Sync()
	akamai.StopSpinnerOk()

	// suppress result output?
	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}

	fmt.Fprintln(c.App.Writer, "")
	fmt.Fprintln(c.App.Writer, results)

	return nil
}

func renderBulkZonesRequestStatusTable(submitStatus *dnsv2.BulkZonesResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Bulk Zones Request Submission Status")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	table.Append([]string{"Request Id", submitStatus.RequestId})
	table.Append([]string{"Expiration Date", submitStatus.ExpirationDate})

	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
