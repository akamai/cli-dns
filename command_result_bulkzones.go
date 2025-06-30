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

func cmdResultBulkZones(c *cli.Context) error {
	ctx := context.Background()

	sess, err := edgegrid.InitializeSession(c)
	if err != nil {
		return fmt.Errorf("session failed %v", err)
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
		return cli.NewExitError(color.RedString("One or more requestids required. "), 1)
	}

	fmt.Println("Preparing bulk zones result request(s) ", "")

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

	var results string
	fmt.Println("Submitting Bulk Zones request  ", "")
	//  Submit
	if op == "create" {
		resultRespCreateList := make([]*dns.GetBulkZoneCreateResultResponse, 0)

		for _, requestid := range requestids {
			resp, err := dnsClient.GetBulkZoneCreateResult(ctx, dns.GetBulkZoneCreateResultRequest{
				RequestID: requestid,
			})
			if err != nil {
				return cli.NewExitError(color.RedString(fmt.Sprintf("bulk zone create error: %s", err)), 1)
			}
			resultRespCreateList = append(resultRespCreateList, resp)
		}
		if c.IsSet("json") && c.Bool("json") {
			jsonData, err := json.MarshalIndent(resultRespCreateList, "", " ")
			if err != nil {
				return cli.NewExitError(color.RedString("Failed to marshal JSON result"), 1)
			}
			results = string(jsonData)
		} else {
			results = renderBulkZonesResultTable(resultRespCreateList, c)
		}
	} else {
		resultRespDeleteList := make([]*dns.GetBulkZoneDeleteResultResponse, 0)
		for _, requestid := range requestids {
			resp, err := dnsClient.GetBulkZoneDeleteResult(ctx, dns.GetBulkZoneDeleteResultRequest{
				RequestID: requestid,
			})
			if err != nil {
				return cli.NewExitError(color.RedString(fmt.Sprintf("bulk zone delete error: %s", err)), 1)
			}
			resultRespDeleteList = append(resultRespDeleteList, resp)
		}
		if c.IsSet("json") && c.Bool("json") {
			jsonData, err := json.MarshalIndent(resultRespDeleteList, "", " ")
			if err != nil {
				return cli.NewExitError(color.RedString("Failed to marshal JSON result"), 1)
			}
			results = string(jsonData)
		} else {
			results = renderBulkZonesResultTable(resultRespDeleteList, c)
		}
	}

	fmt.Println("Assembling Bulk Zone Response Content ", "")
	if len(outputPath) > 1 {
		fmt.Printf("Writing Output to %s ", outputPath)
		// pathname and exists?
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
