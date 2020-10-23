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
	"time"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

const (
	MaxUint     = ^uint(0)
	MaxInt      = int(MaxUint >> 1)
	httpMaxBody = MaxInt
)

func cmdSubmitBulkZones(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	// defer initing dnsv2 until know body size ...

	var (
		outputPath     string
		contractid     string
		groupid        string
		inputPath      string
		bulkDeleteList *dnsv2.ZoneNameListResponse
		newBulkZones   *dnsv2.BulkZonesCreate
		op             string = "create"
		bypass         bool
		maxNumZones    int = 1000
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
			fmt.Println("Warning: bypassZoneSafety arg ignored")
		}
	} else {
		if c.IsSet("contractid") {
			fmt.Println("Warning: contractid arg ignored")
		}
		if c.IsSet("groupid") {
			fmt.Println("Warning: groupid arg ignored")
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

	val, ok := os.LookupEnv("AKAMAI_ZONES_BATCH_SIZE")
	if ok {
		batchsize, err := strconv.Atoi(val)
		if err != nil {
	                akamai.StopSpinnerFail()
        	        return cli.NewExitError(color.RedString(" Environ variable AKAMAI_ZONEBATCH has invalid value"), 1)
        	}
		maxNumZones = batchsize
	}
 
	// Read in json file
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Failed to read input file"), 1)
	}
	// set local variables and Object
	requestMaxBody := config.MaxBody

	if op == "create" {
		err = json.Unmarshal(data, newBulkZones)
	} else {
		err = json.Unmarshal(data, bulkDeleteList)
	}
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Failed to parse json file content into bulk zones object"), 1)
	}
	sdata := string(data)
	if len(sdata) > requestMaxBody {
		requestMaxBody = len(sdata)
	}
	if requestMaxBody > config.MaxBody {
		if requestMaxBody > httpMaxBody {
			config.MaxBody = httpMaxBody
		} else {
			config.MaxBody = requestMaxBody
		}
	}

	// init library with approp sized MaxBody
	dnsv2.Init(config)

	akamai.StartSpinner("Submitting Bulk Zones request  ", "")
	//  Submit
	submitStatusList := make([]*dnsv2.BulkZonesResponse, 0)
	if op == "create" {
		// We should be able to handle up to 1000 zones with even 32bit max body size
		ZonesMax := len(newBulkZones.Zones)
		numZones := ZonesMax
		if ZonesMax > maxNumZones {
			numZones = maxNumZones
		}
		bulkZonesList := make([]*dnsv2.BulkZonesCreate, 0)
		bulkZones := &dnsv2.BulkZonesCreate{Zones: make([]*dnsv2.ZoneCreate, 0)}
		for _, zone := range newBulkZones.Zones {
			bulkZones.Zones = append(bulkZones.Zones, zone)
			numZones -= 1
			ZonesMax -= 1
			if numZones == 0 {
				bulkZonesList = append(bulkZonesList, bulkZones)
				numZones = ZonesMax
				if ZonesMax > maxNumZones {
					numZones = maxNumZones
				}
				bulkZones = &dnsv2.BulkZonesCreate{Zones: make([]*dnsv2.ZoneCreate, 0)}
			}
		}
		for _, zonesRequest := range bulkZonesList {
			submitStatus, err := dnsv2.CreateBulkZones(zonesRequest, queryArgs)
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request submit failed. Error: %s", err.Error())), 1)
			}
			submitStatusList = append(submitStatusList, submitStatus)
		}
	} else {
		// can't imagine size of bulkDeleteList would ever be greater than max body size!
		submitStatus, err := dnsv2.DeleteBulkZones(bulkDeleteList, bypass)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Request submit failed. Error: %s", err.Error())), 1)
		}
		submitStatusList = append(submitStatusList, submitStatus)
	}
	akamai.StopSpinnerOk()
	results := ""
	akamai.StartSpinner("Assembling Bulk Zone Response Content ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(submitStatusList, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display request status"), 1)
		}
		results = string(zjson)
	} else {
		results = renderBulkZonesRequestStatusTable(submitStatusList, c)
	}
	akamai.StopSpinnerOk()

	if len(outputPath) < 1 {
		// if no output path, write out locally. RequestId has to be retreievable
		timeExt := strconv.FormatInt(time.Now().Unix(), 10)
		outputPath = fmt.Sprintf("Bulk_Submit_Request_Status_%s", timeExt)
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

func renderBulkZonesRequestStatusTable(submitStatusList []*dnsv2.BulkZonesResponse, c *cli.Context) string {

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

	for i, submitStatus := range submitStatusList {
		table.Append([]string{"Request Id", submitStatus.RequestId})
		table.Append([]string{"Expiration Date", submitStatus.ExpirationDate})
		if i == len(submitStatusList)-1 {
			table.Append([]string{"", ""})
		}
	}
	table.Render()
	outString += fmt.Sprintln(tableString.String())

	return outString
}
