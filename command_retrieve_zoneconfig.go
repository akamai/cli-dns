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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdRetrieveZoneconfig(c *cli.Context) error {

	fmt.Fprintf(os.Stderr, "Command %s", c.Command.Name)

	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf(color.RedString("session failed %v", err))
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	zonename := c.Args().First()

	// Validate zonename argument
	if zonename == "" {
		return cli.NewExitError(color.RedString("zonename required"), 1)
	}

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

	fmt.Fprintln(os.Stderr, color.BlueString("Retrieving Zone..."))

	zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
	if err != nil {
		if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
			return cli.NewExitError(color.RedString("zone does not exist"), 1)
		}
		return cli.NewExitError(fmt.Sprintf(color.RedString("failed to retrieve zone: %s", err)), 1)
	}

	// Retrieve zone as master zone file
	if isMasterfile {

		// ALIAS zones do not support master file view
		if strings.EqualFold(zone.Type, "ALIAS") {
			return cli.NewExitError(color.RedString(fmt.Sprintf("zone %s is an ALIAS zone and does not support master file retrieval", zonename)), 1)
		}

		content, err := dnsClient.GetMasterZoneFile(ctx, dns.GetMasterZoneFileRequest{Zone: zonename})
		if err != nil {
			if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
				return cli.NewExitError(color.RedString("zone doesn't exist"), 1)
			}
			return cli.NewExitError(fmt.Sprintf(color.RedString("failed to retrive master file: %s", err)), 1)
		}
		results = content
	} else {
		// Retrieve zone in structured format
		/*zone, err := dnsClient.GetZone(ctx, dns.GetZoneRequest{Zone: zonename})
		if err != nil {
			if dnsErr, ok := err.(*dns.Error); ok && dnsErr.StatusCode == 404 {
				return cli.NewExitError(color.RedString("zone does not exist"), 1)
			}
			return cli.NewExitError(fmt.Sprintf(color.RedString("zailed to retrieve zone: %s", err)), 1)
		}*/

		// Output as JSON or table format
		if c.Bool("json") {
			b, err := json.MarshalIndent(zone, "", " ")
			if err != nil {
				return cli.NewExitError(color.RedString("failed to marshal zone JSON"), 1)
			}
			results = string(b)
		} else {
			results = renderZoneconfigTable(zone, c)
		}
	}

	// Write output to file or console
	if outputPath != "" {
		fmt.Fprintf(os.Stderr, color.GreenString("Writing output to %s...\n", outputPath))
		file, err := os.Create(outputPath)
		if err != nil {
			return cli.NewExitError(fmt.Sprintf(color.RedString("failed to create output file: %s", err)), 1)
		}
		defer file.Close()

		if _, err := file.WriteString(results); err != nil {
			return cli.NewExitError(color.RedString("failed to write output to file"), 1)
		}
		return nil
	}

	fmt.Println()
	fmt.Println(results)
	return nil
}
