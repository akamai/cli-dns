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
	"strings"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func cmdResultBulkZones(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		requestids []string
		outputPath string
		op         = "create"
	)

	requestids = c.StringSlice("requestid")
	if len(requestids) < 1 {
		return cli.NewExitError(color.RedString("One or more requestids required. "), 1)
	}

	akamai.StartSpinner("Preparing bulk zones result request(s) ", "")

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
	var resultResp interface{}
	resultRespCreateList := make([]*dnsv2.BulkCreateResultResponse, 0)
	resultRespDeleteList := make([]*dnsv2.BulkDeleteResultResponse, 0)
	for _, requestid := range requestids {
		if op == "create" {
			resultResp, err = dnsv2.GetBulkZoneCreateResult(requestid)
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request Result Query failedd. Error: %s", err.Error())), 1)
			}
			resultRespCreateList = append(resultRespCreateList, resultResp.(*dnsv2.BulkCreateResultResponse))
		} else {
			resultResp, err = dnsv2.GetBulkZoneDeleteResult(requestid)
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request Result Query failedd. Error: %s", err.Error())), 1)
			}
			resultRespDeleteList = append(resultRespDeleteList, resultResp.(*dnsv2.BulkDeleteResultResponse))
		}
	}
	akamai.StopSpinnerOk()

	results := ""
	akamai.StartSpinner("Assembling Bulk Zone Response Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		var zjson []byte
		var err error
		if op == "create" {
			zjson, err = json.MarshalIndent(resultRespCreateList, "", "  ")
		} else {
			zjson, err = json.MarshalIndent(resultRespDeleteList, "", "  ")
		}
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to process result response"), 1)
		}
		results = string(zjson)
	} else {
		if op == "create" {
			results = renderBulkZonesResultTable(resultRespCreateList, c)
		} else {
			results = renderBulkZonesResultTable(resultRespDeleteList, c)
		}
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

func renderBulkZonesResultTable(resultRespList interface{}, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	var requestid string
	var succzones []string
	var failzones []*dnsv2.BulkFailedZone
	op := "Created"
	tableHeader := "Bulk Zones %s Request Results"

	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln(fmt.Sprintf(tableHeader, op))
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	if resultList, ok := resultRespList.([]*dnsv2.BulkCreateResultResponse); ok {
		for _, crreq := range resultList {
			requestid = crreq.RequestId
			succzones = crreq.SuccessfullyCreatedZones
			failzones = crreq.FailedZones
			table.Append([]string{"Request Id", requestid, "", ""})
			table.Append([]string{"", fmt.Sprintf("Successfully %s Zones", op), "", ""})
			if len(succzones) == 0 {
				table.Append([]string{"", "", "None", ""})
			} else {
				for _, zn := range succzones {
					table.Append([]string{"", "", zn, ""})
				}
			}
			table.Append([]string{"", fmt.Sprintf("Failed %s Zones", op), "", ""})
			if len(failzones) == 0 {
				table.Append([]string{"", "", "None", ""})
			} else {
				for _, fzn := range failzones {
					table.Append([]string{"", "", fzn.Zone, fzn.FailureReason})
				}
			}
		}
		table.Render()
		outString += fmt.Sprintln(tableString.String())

		return outString
	}
	resultList, ok := resultRespList.([]*dnsv2.BulkDeleteResultResponse)
	if !ok {
		return "Unable to create result table"
	}
	for _, delreq := range resultList {
		requestid = delreq.RequestId
		succzones = delreq.SuccessfullyDeletedZones
		failzones = delreq.FailedZones
		op = "Deleted"
		table.Append([]string{"Request Id", requestid, "", ""})
		table.Append([]string{fmt.Sprintf("", "Successfully %s Zones", op), "", ""})
		if len(succzones) == 0 {
			table.Append([]string{"", "", "None", ""})
		} else {
			for _, zn := range succzones {
				table.Append([]string{"", "", zn, ""})
			}
		}
		table.Append([]string{fmt.Sprintf("", "Failed %s Zones", op), "", ""})
		if len(succzones) == 0 {
			table.Append([]string{"", "", "None", ""})
		} else {
			for _, fzn := range failzones {
				table.Append([]string{"", "", fzn.Zone, fzn.FailureReason})
			}
		}
	}

	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
