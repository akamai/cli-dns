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
	"github.com/urfave/cli"
)

func cmdUpdateRecordsets(c *cli.Context) error {

	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename   string
		outputPath string
		inputPath  string
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	zonename = c.Args().First()

	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	} else {
		return cli.NewExitError(color.RedString("Input file is required"), 1)
	}
	akamai.StartSpinner("Fetching Recordset data ", "")
	// Read in json file
	data, err := ioutil.ReadFile(inputPath)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Failed to read input file"), 1)
	}
	recordsets := &dnsv2.Recordsets{}
	err = json.Unmarshal(data, recordsets)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString("Failed to parse json file content"), 1)
	}
	// NOTE: UPDATE REPLACES ALL RECORDSETS
	var recordsetWorkList []dnsv2.Recordset
	if c.IsSet("overwrite") && c.Bool("overwrite") {
		recordsets := &dnsv2.Recordsets{}
		err = json.Unmarshal(data, recordsets)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to parse json file content"), 1)
		}
		recordsetWorkList = recordsets.Recordsets
	} else {
		akamai.StartSpinner("Retrieving Existing Recordsets ", "")
		queryArgs := dnsv2.RecordsetQueryArgs{}
		queryArgs.ShowAll = true
		recordsetResp, err := dnsv2.GetRecordsets(zonename, queryArgs)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset List retrieval failed. Error: %s", err.Error())), 1)
		}
		akamai.StopSpinnerOk()
		akamai.StartSpinner("Processing Updated Recordsets ", "")
		recordsetWorkList = recordsetResp.Recordsets
		for _, crs := range recordsets.Recordsets {
			// for each updated recordset
			for i, rs := range recordsetWorkList {
				// walk the full list and relace
				if crs.Name == rs.Name && crs.Type == rs.Type {
					recordsetWorkList[i] = crs
				} else if rs.Type == "SOA" {
					// Serial needs to be incremented
					soavals := strings.Split(rs.Rdata[0], " ")
					v, _ := strconv.Atoi(soavals[2])
					soavals[2] = strconv.Itoa(v + 1)
					rs.Rdata[0] = strings.Join(soavals, " ")
				}
			}
		}
		akamai.StopSpinnerOk()
	}
	akamai.StartSpinner("Updating Recordsets ", "")
	recordsets.Recordsets = recordsetWorkList
	err = recordsets.Update(zonename, true)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset update failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()

	if c.IsSet("suppress") && c.Bool("suppress") {

		return nil

	}
	akamai.StartSpinner("Retrieving Full Recordsets List ", "")
	recordsetResp, err := dnsv2.GetRecordsets(zonename)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset List retrieval failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()
	recordsetList := &RecordsetList{Recordsets: recordsetResp.Recordsets} // list of response objects
	results := ""
	akamai.StartSpinner("Assembling Recordsets List ", "")
	// full output
	if c.IsSet("json") && c.Bool("json") {
		rjson, err := json.MarshalIndent(recordsetList, "", "  ")
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to display recordsets list"), 1)
		}
		results = string(rjson)
	} else {
		results = renderRecordsetListTable(zonename, recordsetList.Recordsets, c)
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
