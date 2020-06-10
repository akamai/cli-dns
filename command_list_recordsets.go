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
	"strconv"
	"strings"
	"path/filepath"
	"os"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

type RecordsetList struct {
	Recordsets	[]dnsv2.Recordset
}

func cmdListRecordsets(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename string
		outputPath string
		rstype []string
		search string
		sortby string = "type"
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	zonename = c.Args().First()
        queryArgs := dnsv2.RecordsetQueryArgs{}
        queryArgs.ShowAll = true

	// for testing
        //queryArgs.ShowAll = false
        //queryArgs.PageSize = 5

	if c.IsSet("sortby") {
		sortby = c.String("sortby")
	}
	queryArgs.SortBy = sortby
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
        if c.IsSet("type") {
                rstype = c.StringSlice("type")
                for i, zt := range rstype {
                        queryArgs.Types += zt
                        if i < len(rstype)-1 {
                                queryArgs.Types += ","
                        }
                }
        }
        if c.IsSet("search") {
                search = c.String("search")
		queryArgs.Search = search
        }
	akamai.StartSpinner("Retrieving Recordsets List ", "")
	recordsetResp, err := dnsv2.GetRecordsets(zonename, queryArgs)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset List retrieval failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()
	recordsets := recordsetResp.Recordsets			// list of response objects
	results := ""
        akamai.StartSpinner("Assembling Recordsets List ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		recordsetList := &RecordsetList{Recordsets: recordsets}
		rjson, err := json.MarshalIndent(recordsetList, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display recordsets list"), 1)
		}
		results = string(rjson)
	} else {
                results = renderRecordsetListTable(zonename, recordsets, c)
	}
	akamai.StopSpinnerOk()
	if len(outputPath) > 1 {
		akamai.StartSpinner(fmt.Sprintf("Writing Output to %s ", outputPath), "")
                rlfHandle, err := os.Create(outputPath)		
                if err != nil {
                        akamai.StopSpinnerFail()
                        return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
                }
		defer rlfHandle.Close()
		_, err = rlfHandle.WriteString(string(results))
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to write zone list output to file"), 1)
		}
		rlfHandle.Sync()
		akamai.StopSpinnerOk()
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}

func renderRecordsetListTable(zone string, recordsets []dnsv2.Recordset, c *cli.Context) string { 

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Recordset List")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetHeader([]string{"NAME", "TYPE", "TTL", "TARGET"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
        table.SetCenterSeparator(" ")
        table.SetColumnSeparator(" ")
        table.SetRowSeparator(" ")
        table.SetBorder(false)
	table.SetCaption(true, fmt.Sprintf("Zone: %s", zone))

	if len(recordsets) == 0 {
		rowData := []string{"No recordsets found", " ", " "}
		table.Append(rowData)
	} else {
		for _, set := range recordsets {
			name := set.Name
			rstype := set.Type
			ttl := strconv.Itoa(set.TTL)
			//rdata := strings.Join(set.Rdata, ", ")
			for i, targ := range set.Rdata {
				if i == 0 {
					table.Append([]string{name, rstype, ttl, targ})
				} else {
                                        table.Append([]string{" ", " ", " ", targ})
				}
			}
			table.Append([]string{" ", " ", " ", " "})
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}

