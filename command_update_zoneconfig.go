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
	"github.com/urfave/cli"
)

func cmdUpdateZoneconfig(c *cli.Context) error {
	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename   string
		outputPath string
		inputPath  string
		masterZoneFileData string
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	akamai.StartSpinner("Preparing zone for update ", "")
	zonename = c.Args().First()
	newZone := &dnsv2.ZoneCreate{}
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	// See if already exists
	zone, err := dnsv2.GetZone(zonename)
	if err != nil {
		if dnsv2.IsConfigDNSError(err) && err.(dnsv2.ConfigDNSError).NotFound() {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Zone already exists"), 1)
		} else {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failure while checking zone existance. Error: %s", err.Error())), 1)
		}
	}
	masterfile :=  c.IsSet("dns") && c.Bool("dns")
	if c.IsSet("file") {
		// Read in json file
		data, err := ioutil.ReadFile(inputPath)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to read input file"), 1)
		}
		if masterfile {
			masterZoneFileData = string(data)
 		} else {
			// set local variables and Object
			err = json.Unmarshal(data, &newZone)
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Failed to parse json file content into zone object"), 1)
			}
			zonename = newZone.Zone
		}
	} else if c.IsSet("type") {
		newZone.Zone = zonename
		newZone.Type = strings.ToUpper(c.String("type"))
		if c.IsSet("contractid") {
			newZone.ContractId = c.String("contractid")
		} else {
			newZone.ContractId = zone.ContractId
		}
		if c.IsSet("master") {
			newZone.Masters = c.StringSlice("master")
		} else {
			newZone.Masters = zone.Masters
		}
		if c.IsSet("comment") {
			newZone.Comment = c.String("comment")
		} else {
			newZone.Comment = zone.Comment
		}
		if c.IsSet("signandserve") {
			newZone.SignAndServe = c.Bool("signandserve")
		} else {
			newZone.SignAndServe = zone.SignAndServe
		}
		if c.IsSet("algorithm") {
			newZone.SignAndServeAlgorithm = c.String("algorithm")
		} else {
			newZone.SignAndServeAlgorithm = zone.SignAndServeAlgorithm
		}
		newZone.TsigKey = zone.TsigKey
		if c.IsSet("tsigname") {
			newZone.TsigKey.Name = c.String("tsigname")
		}
		if c.IsSet("tsigalgorithm") {
			newZone.TsigKey.Algorithm = c.String("tsigalgorithm")
		}
		if c.IsSet("tsigsecret") {
			newZone.TsigKey.Secret = c.String("tsigsecret")
		}
		if c.IsSet("target") {
			newZone.Target = c.String("target")
		} else {
			newZone.Target = zone.Target
		}
		if c.IsSet("endcustomerid") {
			newZone.EndCustomerId = c.String("endcustomerid")
		} else {
			newZone.EndCustomerId = zone.EndCustomerId
		}
	} else {
		akamai.StopSpinnerFail()
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zone command line values or input file are required"), 1)
	}
 
	if masterfile {
	        akamai.StartSpinner("Updating Master Zone File ", "")
		if err = dnsv2.PostMasterZoneFile(zonename, &masterZoneFileData); err != nil {
                        akamai.StopSpinnerFail()
                        return cli.NewExitError(color.RedString(fmt.Sprintf("Master Zone File update failed. Error: %s", err.Error())), 1)
                }
		akamai.StopSpinnerOk()
		return nil
	}

       	akamai.StartSpinner("Updating Zone  ", "")
	err = dnsv2.ValidateZone(newZone)
	if err != nil {
		akamai.StopSpinnerFail()
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString(fmt.Sprintf("Invalid value provided for zone. Error: %s", err.Error())), 1)
	}
	err = newZone.Update(dnsv2.ZoneQueryString{})
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Zone update failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()
	akamai.StartSpinner("Reading Zone Content  ", "")
	zone, err = dnsv2.GetZone(zonename)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read zone content. Error: %s", err.Error())), 1)
	}
        akamai.StopSpinnerOk()

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
