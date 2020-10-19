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
	"strconv"
	"strings"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func cmdStatusBulkZones(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		requestid  string
		outputPath string
		op         = "create"
	)

	if c.IsSet("requestid") {
		requestid = c.String("requestid")
	} else {
		return cli.NewExitError(color.RedString("requestid is required. "), 1)
	}

	akamai.StartSpinner("Preparing bulk zones status request ", "")

	if (c.IsSet("create") && c.IsSet("delete")) || (!c.IsSet("create") && !c.IsSet("delete")) {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Either create or delete arg is required. "), 1)
	}
	if c.IsSet("delete") {
		op = "delete"
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	akamai.StartSpinner("Submitting Bulk Zones request  ", "")
	//  Submit
	var statusResp *dnsv2.BulkStatusResponse
	if op == "create" {
		statusResp, err = dnsv2.GetBulkZoneCreateStatus(requestid)
	} else {
		statusResp, err = dnsv2.GetBulkZoneDeleteStatus(requestid)
	}
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request Status Query failedd. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()

	results := ""
	akamai.StartSpinner("Assembling Bulk Zone Response Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(statusResp, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to process status response"), 1)
		}
		results = string(zjson)
	} else {
		results = renderBulkZonesStatusTable(statusResp, c)
	}
	akamai.StopSpinnerOk()

	if len(outputPath) > 1 {
		akamai.StartSpinner(fmt.Sprintf("Writing Output to %s ", outputPath), "")
		// pathname and exists?
		zfHandle, err := os.Create(outputPath)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer zfHandle.Close()
		_, err = zfHandle.WriteString(string(results))
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to write zone output to file"), 1)
		}
		zfHandle.Sync()
		akamai.StopSpinnerOk()
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil

}

func renderBulkZonesStatusTable(submitStatus *dnsv2.BulkStatusResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Bulk Zones Request Status")
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
	table.Append([]string{"Zones Submitted", strconv.Itoa(submitStatus.ZonesSubmitted)})
	table.Append([]string{"Success Count", strconv.Itoa(submitStatus.SuccessCount)})
	table.Append([]string{"Failure Count", strconv.Itoa(submitStatus.FailureCount)})
	table.Append([]string{"Complete", fmt.Sprintf("%t", submitStatus.IsComplete)})
	table.Append([]string{"Expiration Date", submitStatus.ExpirationDate})
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
