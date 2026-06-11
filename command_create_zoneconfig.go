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
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCreateZoneconfig(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	// Validate zonename argument
	if c.NArg() == 0 {
		return failStep("Preparing zone", "zonename is required")
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()
	zonename := c.Args().First()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing zone", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	// Parse Flags
	var (
		inputPath  = c.String("file")
		outputPath = c.String("output")
		contractID = c.String("contractid")
		groupID    = c.String("groupid")
	)

	newZone := &dns.ZoneCreate{}
	fmt.Printf("Preparing zone ... %s\n", color.GreenString("[OK]"))

	// Load zone config from file if specified
	if inputPath != "" {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return failStep("Preparing zone", "failed to read input file")
		}
		if err := json.Unmarshal(data, newZone); err != nil {
			return failStep("Preparing zone", "failed to parse JSON config")
		}
		//fmt.Printf("Debug: ContractID from JSON: '%s'\n", newZone.ContractID)

		zonename = newZone.Zone

		if contractID == "" {
			contractID = newZone.ContractID
		}
	} else if c.IsSet("type") {
		// Construct zone config from CLI flags
		newZone.Zone = zonename
		newZone.Type = strings.ToUpper(c.String("type"))
		newZone.Comment = c.String("comment")
		newZone.ContractID = c.String("contractid")
		if c.IsSet("master") {
			newZone.Masters = c.StringSlice("master")
		}
		if c.IsSet("signandserve") {
			newZone.SignAndServe = c.Bool("signandserve")
			newZone.SignAndServeAlgorithm = c.String("algorithm")
		}
		if c.IsSet("tsigname") {
			newZone.TSIGKey = &dns.TSIGKey{
				Name:      c.String("tsigname"),
				Algorithm: c.String("tsigalgorithm"),
				Secret:    c.String("tsigsecret"),
			}
		}
		newZone.Target = c.String("target")
		newZone.EndCustomerID = c.String("endcustomerid")

		if contractID == "" {
			contractID = newZone.ContractID
		}
	} else {
		return failStep("Preparing zone", "zone command line values or input file are required")
	}

	if contractID == "" {
		return failStep("Preparing zone", "contractid is required")
	}

	err = dns.ValidateZone(newZone)
	if err != nil {
		return failStep("Preparing zone", "Invalid zone value: %s", err)
	}

	// Check if zone already exists
	_, err = dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err == nil {
		return failStep("Checking Zone Existence", "zone already exists")
	} else {
		if errors.Is(err, dns.ErrGetZone) {
			return failStep("Checking Zone Existence", "failure while checking zone existance")
		}
	}
	fmt.Printf("Checking Zone Existence ... %s\n", color.GreenString("[OK]"))

	// Create new zone
	err = dnsClient.CreateZone(ctx, dns.CreateZoneRequest{
		CreateZone:      newZone,
		ZoneQueryString: dns.ZoneQueryString{Contract: contractID, Group: groupID},
	})
	if err != nil {
		return failStep("Creating Zone", "zone create failed: %s", err)
	}
	fmt.Printf("Creating Zone ... %s\n", color.GreenString("[OK]"))

	// Optionally initialize zone with default records
	if c.Bool("initialize") && strings.ToUpper(newZone.Type) == "PRIMARY" {
		err = dnsClient.SaveChangeList(ctx, dns.SaveChangeListRequest{Zone: zonename})
		if err != nil {
			return failStep("Creating Zone", "failed to initialize zone records")
		}
		err = dnsClient.SubmitChangeList(ctx, dns.SubmitChangeListRequest{Zone: zonename})
		if err != nil {
			return failStep("Creating Zone", "failed to initialize zone records during submit changelist")
		}
	}

	// Fetch zone after creation
	zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		return failStep("Verifying Zone", "failed to read zone config: %v", err)
	}
	fmt.Printf("Verifying Zone ... %s\n", color.GreenString("[OK]"))

	if c.Bool("suppress") {
		return nil
	}

	// Format result for display
	var result string
	if c.Bool("json") {
		b, err := json.MarshalIndent(zone.Zone, "", " ")
		if err != nil {
			return failStep("Assembling Zone Content", "failed to marshal zone output")
		}
		result = string(b)
	} else {
		result = renderZoneconfigTable(zone)
	}
	fmt.Fprintf(os.Stderr, "Assembling Zone Content ... %s\n", color.GreenString("[OK]"))

	// Output to file or stdout
	if outputPath != "" {
		f, err := os.Create(filepath.FromSlash(outputPath))
		if err != nil {
			return failStep("Writing Output", "failed to write output file: %v", err)
		}
		defer func() { _ = f.Close() }()
		_, _ = f.WriteString(result)
		if err := f.Sync(); err != nil {
			return failStep("Writing Output", "failed to sync file: %s", err)
		}
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
	} else {
		_, _ = fmt.Fprintln(c.App.Writer, result)
	}

	return nil
}
