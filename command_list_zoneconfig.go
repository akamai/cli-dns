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

type ZoneSummary struct {
	Zone            string
	Type            string
	ActivationState string
	ContractId      string
}

type ZoneSummaryList struct {
	Zones []*ZoneSummary
}

type ZoneList struct {
	Zones []*dnsv2.ZoneResponse
}

func cmdListZoneconfig(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		outputPath string
		contractid []string
		ztype      []string
		search     string
	)
	queryArgs := dnsv2.ZoneListQueryArgs{}
	queryArgs.ShowAll = true

	// for testing
	//queryArgs.ShowAll = false
	//queryArgs.PageSize = 5

	queryArgs.SortBy = "contractId,type,zone"
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	if c.IsSet("contractid") {
		contractid = c.StringSlice("contractid")
		for i, cid := range contractid {
			queryArgs.ContractIds += cid
			if i < len(contractid)-1 {
				queryArgs.ContractIds += ","
			}
		}
	}
	if c.IsSet("type") {
		ztype = c.StringSlice("type")
		for i, zt := range ztype {
			queryArgs.Types += zt
			if i < len(ztype)-1 {
				queryArgs.Types += ","
			}
		}
	}
	if c.IsSet("search") {
		search = c.String("search")
		queryArgs.Search = search
	}
	akamai.StartSpinner("Retrieving Zone List ", "")
	zoneListResponse, err := dnsv2.ListZones(queryArgs)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Zone List retrieval failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()
	zones := zoneListResponse.Zones // list of ZoneResponse objects
	results := ""
	akamai.StartSpinner("Assembling Zone List ", "")
	if c.IsSet("summary") && c.Bool("summary") {
		if c.IsSet("json") && c.Bool("json") {
			zoneSummary := ZoneSummaryList{}
			zoneSummaryList := make([]*ZoneSummary, 0, len(zones))
			for _, zone := range zones {
				zs := &ZoneSummary{}
				zs.Zone = zone.Zone
				zs.Type = zone.Type
				zs.ActivationState = zone.ActivationState
				zs.ContractId = zone.ContractId
				zoneSummaryList = append(zoneSummaryList, zs)
			}
			zoneSummary.Zones = zoneSummaryList
			json, err := json.MarshalIndent(zoneSummary, "", "  ")
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Unable to display zone list"), 1)
			}
			results = string(json)
		} else {
			results = renderZoneSummaryListTable(zones, c)
		}
	} else {
		// full output
		if c.IsSet("json") && c.Bool("json") {
			zoneSummary := ZoneList{Zones: zones}
			json, err := json.MarshalIndent(zoneSummary, "", "  ")
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Unable to display zone list"), 1)
			}
			results = string(json)
		} else {
			results = renderZoneListTable(zones, c)
		}
	}
	akamai.StopSpinnerOk()
	if len(outputPath) > 1 {
		akamai.StartSpinner(fmt.Sprintf("Writing Output to %s ", outputPath), "")
		zlfHandle, err := os.Create(outputPath)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer zlfHandle.Close()
		_, err = zlfHandle.WriteString(string(results))
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to write zone list output to file"), 1)
		}
		zlfHandle.Sync()
		akamai.StopSpinnerOk()
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}

func renderZoneListTable(zones []*dnsv2.ZoneResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Zone List")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT})
	table.SetHeader([]string{"ZONE", "ATTRIBUTE", "VALUE"})
	table.SetReflowDuringAutoWrap(false)
	table.SetAutoWrapText(false)
	table.SetRowLine(true)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)

	if len(zones) == 0 {
		rowData := []string{"No zones found", " ", " "}
		table.Append(rowData)
	} else {
		for _, zone := range zones {
			zname := zone.Zone
			ztype := zone.Type
			table.Append([]string{zname, "Type", ztype})
			if len(zone.Comment) > 0 {
				table.Append([]string{" ", "Comment", zone.Comment})
			}
			if strings.ToUpper(ztype) == "SECONDARY" {
				if len(zone.Masters) > 0 {
					masters := strings.Join(zone.Masters, " ,")
					table.Append([]string{" ", "Masters", masters})
				}
				if zone.TsigKey != nil {
					table.Append([]string{" ", "TsigKey:Name", zone.TsigKey.Name})
					table.Append([]string{" ", "TsigKey:Algorithm", zone.TsigKey.Algorithm})
					table.Append([]string{" ", "TsigKey:Secret", zone.TsigKey.Secret})
				}
			}
			if strings.ToUpper(ztype) == "PRIMARY" || strings.ToUpper(ztype) == "SECONDARY" {
				table.Append([]string{" ", "SignAndServe", fmt.Sprintf("%t", zone.SignAndServe)})
				if len(zone.SignAndServeAlgorithm) > 0 {
					table.Append([]string{" ", "SignAndServeAlgorithm", fmt.Sprintf("%s", zone.SignAndServeAlgorithm)})
				}
			}
			if strings.ToUpper(ztype) == "ALIAS" {
				table.Append([]string{" ", "Target", zone.Target})
				table.Append([]string{" ", "AliasCount", strconv.FormatInt(zone.AliasCount, 10)})
			}
			table.Append([]string{" ", "ActivationState", zone.ActivationState})
			table.Append([]string{" ", "LastActivationDate", zone.LastActivationDate})
			table.Append([]string{" ", "LastModifiedDate", zone.LastModifiedDate})
			table.Append([]string{" ", "VersionId", zone.VersionId})
			table.Append([]string{" ", " ", " "})
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}

func renderZoneSummaryListTable(zones []*dnsv2.ZoneResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Zone List Summary")
	outString += fmt.Sprintln(" ")
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"ZONE", "TYPE", "ACTIVATION STATE", "CONTRACT ID"})
	table.SetReflowDuringAutoWrap(false)
	table.SetCenterSeparator(" ")
	table.SetColumnSeparator(" ")
	table.SetRowSeparator(" ")
	table.SetBorder(false)
	table.SetAutoWrapText(false)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetRowLine(true)

	if len(zones) == 0 {
		rowData := []string{"No zones found", " ", " ", " "}
		table.Append(rowData)
	} else {
		for _, zone := range zones {
			values := []string{zone.Zone, zone.Type, zone.ActivationState, zone.ContractId}
			table.Append(values)
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString

}
