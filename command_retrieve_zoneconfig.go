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
	"path/filepath"
	"os"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveZoneconfig(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename string
		outputPath string
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

       	akamai.StartSpinner("Preparing zone for create ", "")
	zonename = c.Args().First()
	if c.IsSet("output") {
		outputPath = c.String("output")
                outputPath = filepath.FromSlash(outputPath)
	}

	akamai.StartSpinner("Retrieving Zone  ", "")
	zone, err := dnsv2.GetZone(zonename)
	if err != nil {
		if  dnsv2.IsConfigDNSError(err) && err.(dnsv2.ConfigDNSError).NotFound() {
                	akamai.StopSpinnerFail()
                	return cli.NewExitError(color.RedString("Zone does not exist."), 1)
        	} else {
                	akamai.StopSpinnerFail()
                	return cli.NewExitError(color.RedString(fmt.Sprintf("Zone retrieval failed. Error: %s", err.Error())), 1)
		}

	}	
	akamai.StopSpinnerOk()
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

