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

	"github.com/akamai/cli-dns/edgegrid"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v11/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdStatusBulkZones(c *cli.Context) error {

	// Initialize context Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	var (
		outputPath string
		requestids []string
		op         = "create"
	)

	// Retrieve request IDs from CLI flags
	requestids = c.StringSlice("requestid")
	if len(requestids) < 1 {
		return cli.NewExitError(color.RedString("requestid(s) required. "), 1)
	}

	fmt.Println("Preparing bulk zones status request ", "")

	// Validate that either --create or --delete is set
	if (c.IsSet("create") && c.IsSet("delete")) || (!c.IsSet("create") && !c.IsSet("delete")) {
		return cli.NewExitError(color.RedString("Either create or delete arg is required. "), 1)
	}
	if c.IsSet("delete") {
		op = "delete"
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}

	var statusResp *dns.BulkStatusResponse
	statusRespList := make([]*dns.BulkStatusResponse, 0)
	fmt.Println("Submitting Bulk Zones request(s)  ", "")

	// Loop through all provided request IDs
	for _, requestid := range requestids {
		if op == "create" {
			// Get bulk zone create status
			r, err := dnsClient.GetBulkZoneCreateStatus(ctx, dns.GetBulkZoneCreateStatusRequest{
				RequestID: requestid,
			})
			if err != nil {
				return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Create Status query failed: %s", err)), 1)
			}
			statusResp = &dns.BulkStatusResponse{
				RequestID:      r.RequestID,
				ZonesSubmitted: r.ZonesSubmitted,
				SuccessCount:   r.SuccessCount,
				FailureCount:   r.FailureCount,
				IsComplete:     r.IsComplete,
				ExpirationDate: r.ExpirationDate,
			}
		} else {
			// Get bulk zone delete status
			r, err := dnsClient.GetBulkZoneDeleteStatus(ctx, dns.GetBulkZoneDeleteStatusRequest{
				RequestID: requestid,
			})
			if err != nil {
				return cli.NewExitError(color.RedString(fmt.Sprintf("Bulk Zone Delete Status query failed: %s", err)), 1)
			}
			statusResp = &dns.BulkStatusResponse{
				RequestID:      r.RequestID,
				ZonesSubmitted: r.ZonesSubmitted,
				SuccessCount:   r.SuccessCount,
				FailureCount:   r.FailureCount,
				IsComplete:     r.IsComplete,
				ExpirationDate: r.ExpirationDate,
			}
		}
		statusRespList = append(statusRespList, statusResp)
	}

	results := ""
	fmt.Println("Assembling Bulk Zone Response Content ", "")
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(statusRespList, "", "  ")
		if err != nil {
			return cli.NewExitError(color.RedString("Unable to process status response(s)"), 1)
		}
		results = string(zjson)
	} else {
		results = renderBulkZonesStatusTable(statusRespList, c)
	}

	// Write output to file or console
	if len(outputPath) > 1 {
		//fmt.Printf("Writing Output to %s ", outputPath)
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
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	} else {
		fmt.Fprintln(c.App.Writer, "")
		fmt.Fprintln(c.App.Writer, results)
	}

	return nil

}
