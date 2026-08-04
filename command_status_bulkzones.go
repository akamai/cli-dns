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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/dns"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdStatusBulkZones(c *cli.Context) error {
	// Initialize context Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing bulk zones status request", "session failed %v", err)
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
		return failStep("Preparing bulk zones status request", "requestid(s) required")
	}

	fmt.Printf("Preparing bulk zones status request ... %s\n", color.GreenString("[OK]"))

	// Validate that either --create or --delete is set
	if (c.IsSet("create") && c.IsSet("delete")) || (!c.IsSet("create") && !c.IsSet("delete")) {
		return failStep("Preparing bulk zones status request", "Either create or delete arg is required")
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
	fmt.Printf("Submitting Bulk Zones request(s) ... %s\n", color.GreenString("[OK]"))

	// Loop through all provided request IDs
	for _, requestid := range requestids {
		if op == "create" {
			// Get bulk zone create status
			r, err := dnsClient.GetBulkZoneCreateStatus(ctx, dns.GetBulkZoneCreateStatusRequest{
				RequestID: requestid,
			})
			if err != nil {
				return failStep("Submitting Bulk Zones request(s)", "Bulk Zone Create Status query failed: %s", err)
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
				return failStep("Submitting Bulk Zones request(s)", "Bulk Zone Delete Status query failed: %s", err)
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
	fmt.Printf("Assembling Bulk Zone Response Content ... %s\n", color.GreenString("[OK]"))
	if c.IsSet("json") && c.Bool("json") {
		zjson, err := json.MarshalIndent(statusRespList, "", "  ")
		if err != nil {
			return failStep("Assembling Bulk Zone Response Content", "Unable to process status response(s)")
		}
		results = string(zjson)
	} else {
		results = renderBulkZonesStatusTable(statusRespList)
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
			return failStep("Writing Output", "Unable to write zone output to file")
		}
		_ = zfHandle.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	} else {
		_, _ = fmt.Fprintln(c.App.Writer, results)
	}

	return nil

}
