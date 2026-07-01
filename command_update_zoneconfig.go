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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/akamai/cli-dns/edgegrid"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdUpdateZoneconfig(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing zone", "session failed %v", err)
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
		return failStep("Preparing zone", "zonename is required")
	}

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
		return failStep("Preparing zone", "Either zone command line field values or input file are required")
	}

	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	if c.IsSet("file") {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return failStep("Preparing zone", "Failed to read input file")
		}
		// Update master zone file if dns flag set
		if masterfile {
			masterZoneFileData = string(data)
			if len(masterZoneFileData) > httpMaxBody {
				return failStep("Preparing zone", "Master Zone File size too large to process")
			}
		} else {
			err = json.Unmarshal(data, &newZone)
			if err != nil {
				return failStep("Preparing zone", "Failed to parse json file content into zone object %s", err)
			}
			// Validate required fields from JSON
			if newZone.Zone != "" {
				zonename = strings.TrimSpace(strings.ToLower(newZone.Zone))
			} else {
				return failStep("Preparing zone", "zone is missing in JSON file")
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
		return failStep("Preparing zone", "zone name is required")
	}

	// Fetch zone
	zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		return failStep("Preparing zone", "failure while checking zone existance %s", err)
	}
	if zone == nil {
		return failStep("Preparing zone", "zone retrieval returned nil")
	}
	fmt.Printf("Preparing zone ... %s\n", color.GreenString("[OK]"))

	zoneJson, err := json.MarshalIndent(zone, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal zone for debug: %v\n", err)
	} else {
		fmt.Printf("Retrieved Zone:\n%s\n", string(zoneJson))
	}

	// Handling update using CLI flags
	if c.IsSet("type") && !c.IsSet("file") {
		newZone.Zone = zonename
		newZone.Type = strings.ToUpper(c.String("type"))

		if c.IsSet("contractid") {
			newZone.ContractID = c.String("contractid")
		} else {
			newZone.ContractID = zone.ContractID
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
		if (zone.TSIGKey != nil) || c.IsSet("tsigname") || c.IsSet("tsigalgorithm") || c.IsSet("tsigsecret") {
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
		} else {
			newZone.Target = zone.Target
		}
		if c.IsSet("endcustomerid") {
			newZone.EndCustomerID = c.String("endcustomerid")
		} else {
			newZone.EndCustomerID = zone.EndCustomerID
		}
	}

	// Updating master zone file
	if masterfile {
		err = dnsClient.PostMasterZoneFile(ctx, dns.PostMasterZoneFileRequest{
			Zone:     zonename,
			FileData: masterZoneFileData,
		})
		if err != nil {
			return failStep("Updating Master Zone File", "Master Zone File update failed. Error: %s", err.Error())
		}
		fmt.Printf("Updating Master Zone File ... %s\n", color.GreenString("[OK]"))
		return nil
	}

	//fmt.Printf("DEBUG: updating zone: '%s'\n", newZone.Zone)

	err = dns.ValidateZone(newZone)
	if err != nil {
		return failStep("Updating Zone", "Invalid value provided for zone. Error: %s", err.Error())
	}

	// Updating zone
	err = dnsClient.UpdateZone(ctx, dns.UpdateZoneRequest{
		CreateZone: newZone,
	})
	if err != nil {
		return failStep("Updating Zone", "Zone update failed. Error: %s", err.Error())
	}
	fmt.Printf("Updating Zone ... %s\n", color.GreenString("[OK]"))

	zone, err = dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		return failStep("Verifying Zone", "Failed to read zone content. Error: %s", err.Error())
	}
	fmt.Printf("Verifying Zone ... %s\n", color.GreenString("[OK]"))

	if c.IsSet("suppress") && c.Bool("suppress") {
		return nil
	}
	results := ""

	// Format output either as JSON or table format
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(zone, "", "  ")
		if err != nil {
			return failStep("Assembling Zone Content", "Unable to display zone")
		}
		results = string(zjson)
		fmt.Fprintf(os.Stderr, "Assembling Zone Content ... %s\n", color.GreenString("[OK]"))
	} else {
		fmt.Fprintf(os.Stderr, "Assembling Zone Content ... %s\n", color.GreenString("[OK]"))
		renderZoneTable(zone, nil, c)
		if len(outputPath) > 1 {
			return failStep("Writing Output", "table format output cannot be written to file")
		}
		return nil
	}

	// Write output to file or console
	if len(outputPath) > 1 {
		zfHandle, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "Failed to create output file. Error: %s", err.Error())
		}
		defer func() { _ = zfHandle.Close() }()
		_, err = zfHandle.WriteString(string(results))
		if err != nil {
			return failStep("Writing Output", "Unable to write output to file")
		}
		_ = zfHandle.Sync()
		_, _ = fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	}

	_, _ = fmt.Fprintln(c.App.Writer, results)
	return nil
}
