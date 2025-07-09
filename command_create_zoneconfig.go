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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCreateZoneconfig(c *cli.Context) error {

	// Validate zonename argument
	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	// Initialize context and Edgegrid session
	ctx := context.Background()
	zonename := c.Args().First()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
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

	// Load zone config from file if specified
	if inputPath != "" {
		data, err := os.ReadFile(inputPath)
		if err != nil {
			return cli.NewExitError(color.RedString("failed to read input file"), 1)
		}
		if err := json.Unmarshal(data, newZone); err != nil {
			return cli.NewExitError(color.RedString("failed to parse JSON config"), 1)
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
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zone command line values or input file are required"), 1)
	}

	if contractID == "" {
		return cli.NewExitError(color.RedString("contractid is required"), 1)
	}

	err = dns.ValidateZone(newZone)
	if err != nil {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString(fmt.Sprintf("Invalid zone value: %s", err)), 1)
	}

	// Check if zone already exists
	_, err = dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err == nil {
		return cli.NewExitError(color.RedString("zone already exists"), 1)
	} else {
		if errors.Is(err, dns.ErrGetZone) {
			return cli.NewExitError(color.RedString("failure while checking zone existance"), 1)
		}
	}

	// Create new zone
	err = dnsClient.CreateZone(ctx, dns.CreateZoneRequest{
		CreateZone:      newZone,
		ZoneQueryString: dns.ZoneQueryString{Contract: contractID, Group: groupID},
	})
	if err != nil {
		return cli.NewExitError(color.RedString("zone create failed: %s", err), 1)
	}

	// Optionally initialize zone with default records
	if c.Bool("initialize") && strings.ToUpper(newZone.Type) == "PRIMARY" {
		err = dnsClient.SaveChangeList(ctx, dns.SaveChangeListRequest{Zone: zonename})
		if err != nil {
			return cli.NewExitError(color.RedString("failed to initialize zone records"), 1)
		}
		err = dnsClient.SubmitChangeList(ctx, dns.SubmitChangeListRequest{Zone: zonename})
		if err != nil {
			return cli.NewExitError(color.RedString("failed to initialize zone records during submit changelist "), 1)
		}
	}

	// Fetch zone after creation
	zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("failed to read zone config: %v", err)), 1)
	}

	if c.Bool("suppress") {
		return nil
	}

	// Format result for display
	var result string
	if c.Bool("json") {
		b, err := json.MarshalIndent(zone.Zone, "", " ")
		if err != nil {
			return cli.NewExitError(color.RedString("failed to marshal zone output"), 1)
		}
		result = string(b)
	} else {
		result = renderZoneconfigTable(zone, c)
	}

	// Output to file or stdout
	if outputPath != "" {
		f, err := os.Create(filepath.FromSlash(outputPath))
		if err != nil {
			return cli.NewExitError(color.RedString(fmt.Sprintf("failed to write output file: %v", err)), 1)
		}
		defer f.Close()
		f.WriteString(result)
		f.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, result)
	}

	return nil
}
