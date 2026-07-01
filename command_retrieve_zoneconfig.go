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

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveZoneconfig(c *cli.Context) error {
	failStep := func(step, message string, args ...interface{}) error {
		fmt.Printf("%s ... %s\n", step, color.RedString("[FAIL]"))
		return cli.NewExitError(color.RedString(message, args...), 1)
	}

	//fmt.Fprintf(os.Stderr, "Command %s", c.Command.Name)

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing zone", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	zonename := c.Args().First()

	// Validate zonename argument
	if zonename == "" {
		return failStep("Preparing zone", "zonename required")
	}
	fmt.Printf("Preparing zone ... %s\n", color.GreenString("[OK]"))

	var (
		outputPath string
		results    string
	)

	// Check if the --dns flag is set to retrieve zone as master file
	isMasterfile := c.Bool("dns")

	// Get the output file path if set
	if c.IsSet("output") {
		outputPath = filepath.FromSlash(c.String("output"))
	}

	zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
			return failStep("Retrieving Zone", "zone does not exist")
		}
		return failStep("Retrieving Zone", "failed to retrieve zone: %s", err)
	}
	fmt.Printf("Retrieving Zone ... %s\n", color.GreenString("[OK]"))

	// Retrieve zone as master zone file
	if isMasterfile {

		// ALIAS zones do not support master file view
		if strings.EqualFold(zone.Type, "ALIAS") {
			return failStep("Retrieving Zone", "zone %s is an ALIAS zone and does not support master file retrieval", zonename)
		}

		content, err := dnsClient.GetMasterZoneFile(ctx, dns.GetMasterZoneFileRequest{Zone: zonename})
		if err != nil {
			if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
				return failStep("Retrieving Zone", "zone doesn't exist")
			}
			return failStep("Retrieving Zone", "failed to retrieve master file: %s", err)
		}
		results = content
	} else {
		// Output as JSON or table format
		if c.Bool("json") {
			b, err := json.MarshalIndent(zone, "", " ")
			if err != nil {
				return failStep("Assembling Zone Content", "failed to marshal zone JSON")
			}
			results = string(b)
		} else {
			results = renderZoneconfigTable(zone)
		}
	}
	fmt.Fprintf(os.Stderr, "Assembling Zone Content ... %s\n", color.GreenString("[OK]"))

	// Write output to file or console
	if outputPath != "" {
		//fmt.Fprintf(os.Stderr, color.GreenString("Writing output to %s...\n", outputPath))
		file, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "failed to create output file: %s", err)
		}
		defer func() { _ = file.Close() }()

		if _, err := file.WriteString(results); err != nil {
			return failStep("Writing Output", "failed to write output to file")
		}
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	}

	_, _ = fmt.Fprintln(c.App.Writer, results)
	return nil
}
