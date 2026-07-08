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

func cmdResultBulkZones(c *cli.Context) error {
	// Initialize context and Edgegrid session
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return failStep("Preparing bulk zones result request", "session failed %v", err)
	}
	ctx = edgegrid.WithSession(ctx, sess)
	dnsClient := dns.Client(edgegrid.GetSession(ctx))

	var (
		requestids []string
		outputPath string
		op         = "create"
	)

	requestids = c.StringSlice("requestid")
	if len(requestids) < 1 {
		return failStep("Preparing bulk zones result request", "One or more requestids required")
	}

	// Validate create/delete flags
	if (c.IsSet("create") && c.IsSet("delete")) || (!c.IsSet("create") && !c.IsSet("delete")) {
		return failStep("Preparing bulk zones result request", "Either create or delete arg is required")
	}
	if c.IsSet("delete") {
		op = "delete"
	}
	if c.IsSet("output") {
		outputPath = c.String("output")
		outputPath = filepath.FromSlash(outputPath)
	}
	fmt.Printf("Preparing bulk zones result request ... %s\n", color.GreenString("[OK]"))

	var results string

	if op == "create" {
		resultRespCreateList := make([]*dns.GetBulkZoneCreateResultResponse, 0)
		for _, requestid := range requestids {
			resp, err := dnsClient.GetBulkZoneCreateResult(ctx, dns.GetBulkZoneCreateResultRequest{
				RequestID: requestid,
			})
			if err != nil {
				return failStep("Fetching Bulk Zone Create Results", "bulk zone create error: %s", err)
			}
			resultRespCreateList = append(resultRespCreateList, resp)
		}
		fmt.Printf("Fetching Bulk Zone Create Results ... %s\n", color.GreenString("[OK]"))
		if c.IsSet("json") && c.Bool("json") {
			jsonData, err := json.MarshalIndent(resultRespCreateList, "", " ")
			if err != nil {
				return failStep("Assembling Bulk Zone Response Content", "Failed to marshal JSON result")
			}
			results = string(jsonData)
		} else {
			results = renderBulkZonesResultTable(resultRespCreateList)
		}
	} else {
		resultRespDeleteList := make([]*dns.GetBulkZoneDeleteResultResponse, 0)
		for _, requestid := range requestids {
			resp, err := dnsClient.GetBulkZoneDeleteResult(ctx, dns.GetBulkZoneDeleteResultRequest{
				RequestID: requestid,
			})
			if err != nil {
				return failStep("Fetching Bulk Zone Delete Results", "bulk zone delete error: %s", err)
			}
			resultRespDeleteList = append(resultRespDeleteList, resp)
		}
		fmt.Printf("Fetching Bulk Zone Delete Results ... %s\n", color.GreenString("[OK]"))
		if c.IsSet("json") && c.Bool("json") {
			jsonData, err := json.MarshalIndent(resultRespDeleteList, "", " ")
			if err != nil {
				return failStep("Assembling Bulk Zone Response Content", "Failed to marshal JSON result")
			}
			results = string(jsonData)
		} else {
			results = renderBulkZonesResultTable(resultRespDeleteList)
		}
	}

	fmt.Fprintf(os.Stderr, "Assembling Bulk Zone Response Content ... %s\n", color.GreenString("[OK]"))

	// Write output to file or print to console
	if len(outputPath) > 1 {
		zfHandle, err := os.Create(outputPath)
		if err != nil {
			return failStep("Writing Output", "Failed to create output file. Error: %s", err.Error())
		}
		defer func() { _ = zfHandle.Close() }()
		_, err = zfHandle.WriteString(results)
		if err != nil {
			return failStep("Writing Output", "Unable to write output to file")
		}
		_ = zfHandle.Sync()
		fmt.Fprintln(os.Stderr, color.GreenString("Output written to %s", outputPath))
		return nil
	}

	_, _ = fmt.Fprintln(c.App.Writer, results)
	return nil
}
