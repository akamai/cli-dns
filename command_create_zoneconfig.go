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

func cmdCreateZoneconfig(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename   string
		outputPath string
		contractid string
		groupid    string
		inputPath  string
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	akamai.StartSpinner("Preparing zone for create ", "")
	zonename = c.Args().First()
	newZone := &dnsv2.ZoneCreate{}
	queryArgs := dnsv2.ZoneQueryString{}

	if c.IsSet("contractid") {
		contractid = c.String("contractid")
		newZone.ContractId = contractid
		queryArgs.Contract = contractid
	}
	if c.IsSet("groupid") {
		groupid = c.String("groupid")
		queryArgs.Group = groupid
	}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	if c.IsSet("file") {
		// Read in json file
		data, err := ioutil.ReadFile(inputPath)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to read input file"), 1)
		}
		// set local variables and Object
		err = json.Unmarshal(data, &newZone)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to parse json file content into zone object"), 1)
		}
		zonename = newZone.Zone
		if len(newZone.ContractId) > 0 {
			// overwrite command line arg
			contractid = newZone.ContractId
			queryArgs.Contract = contractid
		}
	} else if c.IsSet("type") {
		// contractid already set
		newZone.Zone = zonename
		newZone.Type = strings.ToUpper(c.String("type"))
		if c.IsSet("master") {
			newZone.Masters = c.StringSlice("master")
		}
		if c.IsSet("comment") {
			newZone.Comment = c.String("comment")
		}
		if c.IsSet("signandserve") {
			newZone.SignAndServe = c.Bool("signandserve")
		}
		if c.IsSet("algorithm") {
			newZone.SignAndServeAlgorithm = c.String("algorithm")
		}
		if c.IsSet("tsigname") {
			newZone.TsigKey = &dnsv2.TSIGKey{}
			newZone.TsigKey.Name = c.String("tsigname")
			if c.IsSet("tsigalgorithm") {
				newZone.TsigKey.Algorithm = c.String("tsigalgorithm")
			}
			if c.IsSet("tsigsecret") {
				newZone.TsigKey.Secret = c.String("tsigsecret")
			}
		}
		if c.IsSet("target") {
			newZone.Target = c.String("target")
		}
		if c.IsSet("endcustomerid") {
			newZone.EndCustomerId = c.String("endcustomerid")
		}
	} else {
		akamai.StopSpinnerFail()
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zone command line values or input file are required"), 1)
	}
	if len(contractid) == 0 {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("contractid is required"), 1)
	}
	err = dnsv2.ValidateZone(newZone)
	if err != nil {
		akamai.StopSpinnerFail()
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString(fmt.Sprintf("Invalid value provided for zone. Error: %s", err.Error())), 1)
	}

	akamai.StartSpinner("Creating Zone  ", "")
	// See if already exists
	zone, err := dnsv2.GetZone(zonename)
	if err == nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Zone already exists"), 1)
	} else {
		if !dnsv2.IsConfigDNSError(err) || !err.(dnsv2.ConfigDNSError).NotFound() {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failure while checking zone existance. Error: %s", err.Error())), 1)
		}
	}
	// create
	err = newZone.Save(queryArgs)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Zone create failed. Error: %s", err.Error())), 1)
	}
	if c.IsSet("initialize") && c.Bool("initialize") && strings.ToUpper(newZone.Type) == "PRIMARY" {
		// Indirectly create NS and SOA records
		err = newZone.SaveChangelist()
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Zone initialization failed. SOA and NS records need to be created "), 1)
		}
		err = newZone.SubmitChangelist()
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Zone create failed. Error: %s", err.Error())), 1)
		}
	}
	akamai.StopSpinnerOk()
	akamai.StartSpinner("Reading Zone Content  ", "")
	zone, err = dnsv2.GetZone(zonename)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read zone content. Error: %s", err.Error())), 1)
	}
	// suppress result output?
	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	results := ""
	akamai.StartSpinner("Assembling Zone Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display zone"), 1)
		}
		results = string(zjson)
	} else {
		results = renderZoneconfigTable(zone, c)
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

func renderZoneconfigTable(zone *dnsv2.ZoneResponse, c *cli.Context) string {

	//bold := color.New(color.FgWhite, color.Bold)
	outString := ""
	outString += fmt.Sprintln(" ")
	outString += fmt.Sprintln("Zone Configuration")
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

	if zone == nil {
		rowData := []string{"No zone info to display", " ", " "}
		table.Append(rowData)
	} else {
		zname := zone.Zone
		ztype := zone.Type
		table.Append([]string{zname, "Type", ztype})
		if len(zone.Comment) > 0 {
			table.Append([]string{" ", "Comment", zone.Comment})
		}
		if len(zone.ContractId) > 0 {
			table.Append([]string{" ", "ContractId", zone.ContractId})
		}
		if strings.ToUpper(ztype) == "SECONDARY" {
			if len(zone.Masters) > 0 {
				masters := strings.Join(zone.Masters, " ,")
				table.Append([]string{" ", "Masters", masters})
			}
			if zone.TsigKey != nil {
				if len(zone.TsigKey.Name) > 0 {
					table.Append([]string{" ", "TsigKey:Name", zone.TsigKey.Name})
				}
				if len(zone.TsigKey.Algorithm) > 0 {
					table.Append([]string{" ", "TsigKey:Algorithm", zone.TsigKey.Algorithm})
				}
				if len(zone.TsigKey.Secret) > 0 {
					table.Append([]string{" ", "TsigKey:Secret", zone.TsigKey.Secret})
				}
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
		if len(zone.LastActivationDate) > 0 {
			table.Append([]string{" ", "LastActivationDate", zone.LastActivationDate})
		}
		if len(zone.LastModifiedDate) > 0 {
			table.Append([]string{" ", "LastModifiedDate", zone.LastModifiedDate})
		}
		table.Append([]string{" ", "VersionId", zone.VersionId})
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
