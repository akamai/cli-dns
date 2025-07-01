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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateZoneconfig(c *cli.Context) error {

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	var (
		zonename           string
		outputPath         string
		inputPath          string
		masterZoneFileData string
	)

	// Validate zonename argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	fmt.Println("Preparing zone for update ", "")

	zonename = c.Args().First()

	// New zone struct to hold zone data for update
	newZone := &dns.ZoneCreate{}

	masterfile := c.IsSet("dns") && c.Bool("dns")

	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
		if c.IsSet("type") {
			fmt.Println("Warning: Zone Field and File args are defined. Field values will be ignored!")
		}
	} else if !c.IsSet("type") && !masterfile {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("Either zone command line field values or input file are required"), 1)
	}

	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	if c.IsSet("file") {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return cli.NewExitError(color.RedString("Failed to read input file"), 1)
		}
		// Update master zone file if dns flag set
		if masterfile {
			masterZoneFileData = string(data)
			if len(masterZoneFileData) > httpMaxBody {
				return cli.NewExitError(color.RedString("Master Zone File size too large to process"), 1)
			}
		} else {
			err = json.Unmarshal(data, &newZone)
			if err != nil {
				return cli.NewExitError(color.RedString("Failed to parse json file content into zone object %s", err), 1)
			}
			// Validate required fields from JSON
			if newZone.Zone != "" {
				zonename = strings.TrimSpace(strings.ToLower(newZone.Zone))
			} else {
				return cli.NewExitError(color.RedString("zone is missing in JSON file"), 1)
			}
			if newZone.Type != "" {
				newZone.Type = strings.ToUpper(newZone.Type)
			}
			if newZone.SignAndServeAlgorithm != "" {
				newZone.SignAndServeAlgorithm = strings.ToUpper(newZone.SignAndServeAlgorithm) // Uppercase signAndServeAlgorithm
			}
		}
	}

	if zonename == "" {
		return cli.NewExitError(color.RedString("zone name is required"), 1)
	}

	// Fetch zone
	zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		return cli.NewExitError(color.RedString("failure while checking zone existance %s", err), 1)
	}
	if zone == nil {
		return cli.NewExitError(color.RedString("zone retrieval returned nil!"), 1)
	}

	zoneJson, err := json.MarshalIndent(zone, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal zone for debug: %v\n", err)
	} else {
		fmt.Printf("Retrieved Zone:\n%s\n", string(zoneJson))
	}

	payload, _ := json.MarshalIndent(newZone, "", "  ")
	fmt.Println("Payload to be sent:\n", string(payload))

	// Handling update using CLI flags
	if c.IsSet("type") && !c.IsSet("file") {
		newZone.Zone = zonename
		newZone.Type = strings.ToUpper(c.String("type"))

		if c.IsSet("contractid") {
			newZone.ContractID = c.String("contractid")
		} else if zone != nil {
			newZone.ContractID = zone.ContractID
		}
		if c.IsSet("master") {
			newZone.Masters = c.StringSlice("master")
		} else if zone != nil {
			newZone.Masters = zone.Masters
		}
		if c.IsSet("comment") {
			newZone.Comment = c.String("comment")
		} else if zone != nil {
			newZone.Comment = zone.Comment
		}
		if c.IsSet("signandserve") {
			newZone.SignAndServe = c.Bool("signandserve")
		} else if zone != nil {
			newZone.SignAndServe = zone.SignAndServe
		}
		if c.IsSet("algorithm") {
			newZone.SignAndServeAlgorithm = c.String("algorithm")
		} else if zone != nil {
			newZone.SignAndServeAlgorithm = zone.SignAndServeAlgorithm
		}
		if (zone != nil && zone.TSIGKey != nil) || c.IsSet("tsigname") || c.IsSet("tsigalgorithm") || c.IsSet("tsigsecret") {
			if zone.TSIGKey != nil {
				newZone.TSIGKey = &dns.TSIGKey{
					Name:      zone.TSIGKey.Name,
					Algorithm: zone.TSIGKey.Algorithm,
					Secret:    zone.TSIGKey.Secret,
				}
			} else {
				newZone.TSIGKey = &dns.TSIGKey{}
			}

			if c.IsSet("tsigname") {
				newZone.TSIGKey.Name = c.String("tsigname")
			}
			if c.IsSet("tsigalgorithm") {
				newZone.TSIGKey.Algorithm = c.String("tsigalgorithm")
			}
			if c.IsSet("tsigsecret") {
				newZone.TSIGKey.Secret = c.String("tsigsecret")
			}
		}
		if c.IsSet("target") {
			newZone.Target = c.String("target")
		} else if zone != nil {
			newZone.Target = zone.Target
		}
		if c.IsSet("endcustomerid") {
			newZone.EndCustomerID = c.String("endcustomerid")
		} else if zone != nil {
			newZone.EndCustomerID = zone.EndCustomerID
		}
	}

	// Updating master zone file
	if masterfile {
		fmt.Println("Updating Master Zone File ", "")
		err = dnsClient.PostMasterZoneFile(ctx, dns.PostMasterZoneFileRequest{
			Zone:     zonename,
			FileData: masterZoneFileData,
		})
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Master Zone File update failed. Error: %s", err.Error())), 1)
		}
		return nil
	}

	fmt.Printf("DEBUG: updating zone: '%s'\n", newZone.Zone)

	fmt.Println("Updating Zone  ", "")
	err = dns.ValidateZone(newZone)

	if err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString(fmt.Sprintf("Invalid value provided for zone. Error: %s", err.Error())), 1)
	}

	// Updating zone
	err = dnsClient.UpdateZone(ctx, dns.UpdateZoneRequest{
		CreateZone: newZone,
	})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Zone update failed. Error: %s", err.Error())), 1)
	}

	fmt.Println("Reading Zone Content  ", "")
	zone, err = dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to read zone content. Error: %s", err.Error())), 1)
	}

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	results := ""
	fmt.Println("Assembling Zone Content ", "")

	// Format output either as JSON or table format
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to display zone"), 1)
		}
		results = string(zjson)
	} else {
		results = renderZoneconfigTable(zone, c)
	}

	// Write output to file or console
	if len(outputPath) > 1 {
		fmt.Printf("Writing Output to %s ", outputPath)
		zfHandle, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("Failed to create output file. Error: %s", err.Error())), 1)
		}
		defer zfHandle.Close()
		_, err = zfHandle.WriteString(string(results))
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to write zone output to file"), 1)
		}
		zfHandle.Sync()
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil
}
